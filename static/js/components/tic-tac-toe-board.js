export class TicTacToeBoard extends HTMLElement {
  size = 3
  board = Array(this.size * this.size).fill(null)
  wonPosition = null
  player = 'X'
  playerTurn = 'X'

  constructor () {
    super()
    this.attachShadow({ mode: 'open' })
    this.shadowRoot.innerHTML = `
      <style>
        .board {
          display: grid;
          grid-template-columns: repeat(var(--size, 3), 1fr);
          grid-gap: calc(30px / var(--size, 3));
          width: calc(var(--size, 3) * 100px);
          height: calc(var(--size, 3) * 100px);
          max-width: 600px;
          max-height: 600px;
          margin: 0 auto;
        }

        .cell {
          position: relative;
          color: #111;
          background-color: #f0f0f0;
        }
        .cell:not(.cell--filled):hover {
          cursor: pointer;
          background-color: #e0e0e0;
        }
        
        .cell--won {
          color: #e55;
        }
        
        .cell::before {
          content: '';
          position: absolute;
          top: 50%;
          left: 50%;
          opacity: 0;
          transform: scale(.8) translate(-50%, -50%);
          transition: .1s;
        }
        .cell--x::before {
          opacity: 1;
          transform: scale(1) translate(-50%, -50%) rotate(45deg);
          width: 80%;
          height: 80%;
          background: 
            linear-gradient(to bottom, transparent 46%, currentColor 46%, currentColor 54%, transparent 54%),
            linear-gradient(to left, transparent 46%, currentColor 46%, currentColor 54%, transparent 54%);
        }
        .cell--o::before {
          opacity: 1;
          transform: scale(1) translate(-50%, -50%);
          width: 80%;
          height: 80%;
          background: radial-gradient(circle at 50% 50%, transparent 45%, currentColor 45%, currentColor 55%, transparent 55%);
        }
      </style>
      <div class="board"></div>
    `

    this.setSize(3)
  }

  get boardEl () {
    return this.shadowRoot.querySelector('.board')
  }

  get cellEls () {
    const cells = []

    this.boardEl.childNodes.forEach((cell) => {
      cells[cell.dataset.index] = cell
    })

    return cells
  }

  setSize (size) {
    const newBoard = new Array(size * size).fill(null)

    const xOffSet = Math.floor((size - this.size) / 2)
    const yOffSet = Math.floor((size - this.size) / 2)
    for (let x = 0; x < this.size; x++) {
      for (let y = 0; y < this.size; y++) {
        newBoard[(x + xOffSet) + (y + yOffSet) * size] = this.board[x + y * this.size] || null
      }
    }

    this.size = size
    this.board = newBoard
    this.wonPosition = null
    this.style.setProperty('--size', size)
    this.renderCells()
  }

  setValue (index, value) {
    this.board[index] = value
    this.updateCell(index)

    const win = this.checkWin()
    if (win) {
      this.wonPosition = win
      win.map((i) => this.updateCell(i))
    }

    this.dispatchEvent(new CustomEvent('turn', { detail: { index } }))

    this.playerTurn = value === 'X' ? 'O' : 'X'
  }

  renderCells () {
    this.boardEl.innerHTML = ''
    for (let i = 0; i < this.size * this.size; i++) {
      this.renderCell(i)
    }
  }

  renderCell (index) {
    const cell = document.createElement('div')

    cell.classList.add('cell')
    cell.dataset.index = index
    this.boardEl.appendChild(cell)
    this.updateCell(index)

    cell.addEventListener('click', () => this.onCellClick(index))
  }

  updateCell (index) {
    const cell = this.cellEls[index]
    cell.classList.add('cell')

    if (this.wonPosition && this.wonPosition.includes(index)) {
      cell.classList.add('cell--won')
    }

    if (this.board[index]) {
      cell.classList.add('cell--' + this.board[index].toLowerCase())
      cell.classList.add('cell--filled')
    } else {
      cell.classList.remove('cell--filled')
      cell.classList.remove('cell--x')
      cell.classList.remove('cell--o')
    }
  }

  onCellClick (index) {
    if (this.board[index] || this.wonPosition || this.playerTurn !== this.player) {
      return
    }

    this.setValue(index, this.player)
  }

  checkWin () {
    const winLength = this.size <= 4 ? this.size : this.size - 1
    const lines = this.getWinPositions(this.size, winLength)

    for (const line of lines) {
      const symbols = line.map((index) => this.board[index])
      const targetSymbol = symbols[0]

      if (['X', 'O'].includes(targetSymbol) && symbols.every((symbol) => symbol === targetSymbol)) {
        return line
      }
    }

    return false
  }

  getWinPositions (size, winLength) {
    const res = []

    // Vertical
    for (let xOffset = 0; xOffset <= size - winLength; xOffset++) {
      for (let x = 0; x < winLength; x++) {
        const column = []

        for (let y = 0; y < winLength; y++) {
          column.push((y * size) + (x + xOffset))
        }

        res.push(column)
      }
    }

    // Horizontal
    for (let yOffset = 0; yOffset <= size - winLength; yOffset++) {
      for (let y = 0; y < winLength; y++) {
        const row = []

        for (let x = 0; x < winLength; x++) {
          row.push((y + yOffset) * size + x)
        }

        res.push(row)
      }
    }

    // Diagonal
    for (let xOffset = 0; xOffset <= size - winLength; xOffset++) {
      for (let yOffset = 0; yOffset <= size - winLength; yOffset++) {
        const diagonal1 = []
        const diagonal2 = []

        for (let i = 0; i < winLength; i++) {
          diagonal1.push(xOffset + i + (yOffset + i) * size)
          diagonal2.push(xOffset + winLength - 1 - i + (yOffset + i) * size)
        }

        res.push(diagonal1)
        res.push(diagonal2)
      }
    }

    return res
  }
}
