# Halooid UI Components

A shared UI component library for the Halooid platform.

## Installation

```bash
npm install @halooid/ui-components
```

## Usage

```jsx
import { Button, Card, Input } from '@halooid/ui-components';
import { ThemeProvider } from 'styled-components';
import { theme } from '@halooid/ui-components';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <Card elevation="md" hoverable>
        <Card.Header>
          <h3 className="card-title">Login</h3>
        </Card.Header>
        <Card.Content>
          <Input 
            label="Email" 
            placeholder="Enter your email" 
            fullWidth 
          />
          <Input 
            label="Password" 
            type="password" 
            placeholder="Enter your password" 
            fullWidth 
            style={{ marginTop: '1rem' }}
          />
        </Card.Content>
        <Card.Footer>
          <Button variant="primary" fullWidth>Login</Button>
        </Card.Footer>
      </Card>
    </ThemeProvider>
  );
}
```

## Components

### Button

A customizable button component.

```jsx
<Button 
  variant="primary" 
  size="md" 
  onClick={() => console.log('Button clicked')}
>
  Click Me
</Button>
```

### Input

A text input component with various styles and states.

```jsx
<Input 
  label="Email" 
  placeholder="Enter your email" 
  helperText="We'll never share your email with anyone else." 
/>
```

### Card

A container component for displaying content in a contained format.

```jsx
<Card elevation="md" hoverable>
  <Card.Header>
    <h3 className="card-title">Card Title</h3>
  </Card.Header>
  <Card.Content>
    <p>This is the card content.</p>
  </Card.Content>
  <Card.Footer>
    <Button variant="primary">Action</Button>
  </Card.Footer>
</Card>
```

## Theme

The component library comes with a default theme that can be customized.

```jsx
import { ThemeProvider } from 'styled-components';
import { theme } from '@halooid/ui-components';

// Customize the theme
const customTheme = {
  ...theme,
  colors: {
    ...theme.colors,
    primary: {
      main: '#ff0000',
      light: '#ff3333',
      dark: '#cc0000',
      contrastText: '#ffffff',
    },
  },
};

function App() {
  return (
    <ThemeProvider theme={customTheme}>
      {/* Your app components */}
    </ThemeProvider>
  );
}
```

## Development

### Setup

```bash
# Install dependencies
npm install

# Start Storybook
npm run storybook
```

### Building

```bash
# Build the component library
npm run build
```

### Testing

```bash
# Run tests
npm test
```

## License

MIT
