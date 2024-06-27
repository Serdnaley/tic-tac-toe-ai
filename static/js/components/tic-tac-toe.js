const statusEnum = {
  WAITING_FOR_MAP_BUILD: 'WAITING_FOR_MAP_BUILD',
  WAITING_FOR_PLAYER_TURN: 'WAITING_FOR_PLAYER_TURN',
  WAITING_FOR_NEXT_TURN_FROM_SERVER: 'WAITING_FOR_NEXT_TURN_FROM_SERVER',
  PLAYER_WON: 'PLAYER_WON',
  PLAYER_LOST: 'PLAYER_LOST',
  DRAW: 'DRAW',
  CRASHED: 'CRASHED',
}

export class TicTacToe extends HTMLElement {
  #status = statusEnum.WAITING_FOR_PLAYER_TURN
  #mapBuildProgress = 0

  constructor () {
    super()
    this.attachShadow({ mode: 'open' })
    this.shadowRoot.innerHTML = `
      <style>
        tic-tac-toe-board {
          width: 100%;
          height: 100%; 
        }
        button {
          margin-left: 5px;
          background: #f0f0f0;
          border: none;
          padding: 5px 10px;
        }
        button:hover {
          background: #e0e0e0;
        }
      </style>
      <p>
        <span>Server:</span>
        <span class="status">Loading...</span>
      </p>
      <p>
        <span>Your chances:</span>
        <span class="chances"></span>
      </p>
      <tic-tac-toe-board></tic-tac-toe-board>
    `

    this.setStatus(statusEnum.WAITING_FOR_PLAYER_TURN)

    this.board.addEventListener('turn', (event) => {
      const { index } = event.detail
      const player = this.board.board[index]

      if (this.board.wonPosition) {
        return this.setStatus(player === 'X' ? statusEnum.PLAYER_WON : statusEnum.PLAYER_LOST)
      }

      if (!this.board.board.some(i => !i)) {
        return this.setStatus(statusEnum.DRAW)
      }

      if (player === 'X') {
        return this.nextMove(this.board.board).catch((err) => this.onError(err))
      }
    })
  }

  get board () {
    return this.shadowRoot.querySelector('tic-tac-toe-board')
  }

  setStatus (status) {
    this.#status = status

    const statusEl = this.shadowRoot.querySelector('.status')
    statusEl.innerHTML = {
      [statusEnum.WAITING_FOR_MAP_BUILD]: `Building the map... (${this.#mapBuildProgress.toFixed(0)}%)`,
      [statusEnum.WAITING_FOR_PLAYER_TURN]: 'Waiting for your move...',
      [statusEnum.WAITING_FOR_NEXT_TURN_FROM_SERVER]: 'Thinking...',
      [statusEnum.PLAYER_WON]: 'You won, but you are a cheater!',
      [statusEnum.PLAYER_LOST]: 'You are a loooooser!',
      [statusEnum.DRAW]: 'Let\'s call it a draw.',
      [statusEnum.CRASHED]: 'The server crashed. Please refresh the page.',
    }[status] || status

    if ([
      statusEnum.PLAYER_WON,
      statusEnum.PLAYER_LOST,
      statusEnum.DRAW,
    ].includes(status)) {
      const btn = document.createElement('button')
      btn.textContent = 'Restart'
      btn.onclick = () => window.location.reload()
      statusEl.appendChild(btn)
    }

    if (statusEnum.DRAW === status) {
      const btn = document.createElement('button')
      btn.textContent = 'Continue'
      btn.onclick = () => {
        this.board.setSize(this.board.size + 2)
        if (this.board.playerTurn === 'O') {
          this.nextMove().catch((err) => this.onError(err))
        }
      }
      statusEl.appendChild(btn)
    }
  }

  getGameString () {
    return [
      this.board.wonPosition?.[0] || '_',
      this.board.board.map(i => i || '_').join(''),
    ].join(' ')
  }

  async waitForMapBuild () {
    this.setStatus(statusEnum.WAITING_FOR_MAP_BUILD)

    const statusRes = await fetch(`/api/maps/status?game=${this.getGameString()}`, { method: 'GET' })
      .then((res) => res.json())
      .catch((err) => this.onError(err))

    if (this.#status !== statusEnum.WAITING_FOR_MAP_BUILD) {
      return
    }

    if (statusRes.data.progress === 100) {
      return this.nextMove().catch((err) => this.onError(err))
    }

    this.#mapBuildProgress = statusRes.data.progress

    this.setStatus(statusEnum.WAITING_FOR_MAP_BUILD)
    setTimeout(() => {
      this.waitForMapBuild().catch((err) => this.onError(err))
    }, 1000)
  }

  async nextMove () {
    this.setStatus(statusEnum.WAITING_FOR_NEXT_TURN_FROM_SERVER)

    const statusRes = await fetch(`/api/maps/status?game=${this.getGameString()}`, { method: 'GET' })
      .then((res) => res.json())
      .catch((err) => this.onError(err))

    if (statusRes.data.progress === 0) {
      await fetch(`/api/maps/build?game=${this.getGameString()}`, { method: 'POST' })
        .catch((err) => this.onError(err))
      return this.waitForMapBuild().catch((err) => this.onError(err))
    }

    if (statusRes.data.progress < 100) {
      return this.waitForMapBuild().catch((err) => this.onError(err))
    }

    const moveRes = await fetch(`/api/next-move?game=${this.getGameString()}`, { method: 'GET' })
      .then((res) => res.json())
      .catch((err) => this.onError(err))
    const { x, y } = moveRes.data
    const i = x + y * this.board.size

    if (this.#status !== statusEnum.WAITING_FOR_NEXT_TURN_FROM_SERVER) {
      return
    }

    if (this.board.board[i]) {
      console.error('Invalid move received from the server.')
      return this.setStatus(statusEnum.CRASHED)
    }
    this.board.setValue(i, 'O')
    this.updateChances().catch((err) => this.onError(err))

    if (this.board.playerTurn === 'X' && !this.board.wonPosition) {
      return this.setStatus(statusEnum.WAITING_FOR_PLAYER_TURN)
    }
  }

  async updateChances () {
    this.shadowRoot.querySelector('.chances').textContent = 'Loading...'

    const chancesRes = await fetch(`/api/chances?game=${this.getGameString()}`, { method: 'GET' })
      .then((res) => res.json())
      .catch((err) => this.onError(err))

    if (this.#status === statusEnum.CRASHED) {
      return
    }

    const { win, lose, draw } = chancesRes.data
    const total = win + lose + draw
    const winPercent = win / total * 100
    const losePercent = lose / total * 100
    const drawPercent = draw / total * 100

    this.shadowRoot.querySelector('.chances').textContent = [
      `Win: ${winPercent.toFixed(0)}%`,
      `Lose: ${losePercent.toFixed(0)}%`,
      `Draw: ${drawPercent.toFixed(0)}%`,
    ].join(', ')
  }

  onError (error) {
    console.error(error)
    this.setStatus(statusEnum.CRASHED)
  }
}