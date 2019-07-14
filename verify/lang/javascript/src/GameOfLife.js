// Size and Compute length parameters
const rows = 5
const cols = 10
const iterations = 20

// not [x][y] [y][x]
const arra = []
const arrb = []
for (let i = 0; i < rows; i++) {
  arra[i] = []
  arrb[i] = []
  for (let j = 0; j < cols; j++) {
    arra[i][j] = 0
  }
}

// Add the example Glider
arra[1][3] = 1
arra[2][1] = 1
arra[2][3] = 1
arra[3][2] = 1
arra[3][3] = 1

// Print Function
function print (array) {
  for (let y = 0; y < rows; y++) {
    let mess = ''
    for (let x = 0; x < cols; x++) {
      mess += array[y][x]
    }
    console.log(mess)
  }
  console.log()
}

// Cell Check Logic
function cell (y, x) {
  let r = arrb[y][x]
  let c = -r
  // Get the amount of cells around this one
  for (let i = -1; i <= 1; i++) {
    for (let j = -1; j <= 1; j++) {
      const xj = x + j + cols
      const cx = xj % cols
      const yi = y + i + rows
      const cy = yi % rows
      c += arrb[cy][cx]
    }
  }
  // Decide the cells fate based on that
  if (c < 2 || c > 3) {
    r = 0
  }
  if (c === 3) {
    r = 1
  }
  return r
}

// Mainloop
print(arra)
for (let g = 0; g <= iterations; g++) {
  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < cols; x++) {
      arrb[y][x] = arra[y][x]
    }
  }
  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < cols; x++) {
      arra[y][x] = cell(y, x)
    }
  }
}
print(arra)
