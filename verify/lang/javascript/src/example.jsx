const MenuItem = require('MenuItem')
const Menu = require('Menu')
const Dropdown = require('Dropdown')

// Using JSX to express UI components.
const dropdown = (
  <Dropdown>
    A dropdown list
    <Menu>
      <MenuItem>Do Something</MenuItem>
      <MenuItem>Do Something Fun!</MenuItem>
      <MenuItem>Do Something Else</MenuItem>
    </Menu>
  </Dropdown>
)

const z = { foo: 'bar' }

console.log({ z, dropdown })
