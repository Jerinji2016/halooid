# Web UI Component Library

## Overview

The Halooid platform includes a shared UI component library for web applications. This library provides a consistent look and feel across all web interfaces in the platform.

## Technology Stack

- **React**: UI library
- **TypeScript**: Type-safe JavaScript
- **Styled Components**: CSS-in-JS styling
- **Storybook**: Component documentation and testing
- **Rollup**: Module bundling

## Installation

```bash
# From the project root
cd web/packages/ui-components

# Install dependencies
npm install

# Build the library
npm run build

# Run Storybook
npm run storybook
```

## Components

### Button

A versatile button component with various styles and states.

```jsx
import { Button } from '@halooid/ui-components';

// Primary button
<Button variant="primary" onClick={() => console.log('Clicked!')}>
  Click Me
</Button>

// Secondary button
<Button variant="secondary" outlined>
  Secondary Button
</Button>

// Disabled button
<Button variant="primary" disabled>
  Disabled Button
</Button>

// Loading button
<Button variant="primary" loading>
  Loading...
</Button>

// Button with icon
<Button 
  variant="primary" 
  startIcon={<Icon name="plus" />}
>
  Add Item
</Button>
```

#### Props

- `variant`: Button style variant (`primary`, `secondary`, `success`, `danger`, `warning`, `info`)
- `size`: Button size (`sm`, `md`, `lg`)
- `outlined`: Whether the button is outlined
- `rounded`: Whether the button is rounded
- `disabled`: Whether the button is disabled
- `fullWidth`: Whether the button is full width
- `loading`: Whether the button is in a loading state
- `startIcon`: Icon to display before the button text
- `endIcon`: Icon to display after the button text
- `onClick`: Click handler function

### Input

A customizable text input component with various styles and states.

```jsx
import { Input } from '@halooid/ui-components';

// Basic input
<Input 
  label="Email" 
  placeholder="Enter your email"
/>

// Input with error
<Input 
  label="Email" 
  placeholder="Enter your email"
  error="Please enter a valid email address"
/>

// Input with helper text
<Input 
  label="Email" 
  placeholder="Enter your email"
  helperText="We'll never share your email with anyone else"
/>

// Input with icon
<Input 
  label="Email" 
  placeholder="Enter your email"
  startIcon={<Icon name="mail" />}
/>
```

#### Props

- `label`: Input label
- `placeholder`: Input placeholder
- `value`: Input value
- `onChange`: Change handler function
- `error`: Error message
- `helperText`: Helper text
- `disabled`: Whether the input is disabled
- `fullWidth`: Whether the input is full width
- `startIcon`: Icon to display at the start of the input
- `endIcon`: Icon to display at the end of the input
- `type`: Input type (`text`, `password`, `email`, etc.)
- `variant`: Input variant (`outlined`, `filled`, `standard`)
- `size`: Input size (`sm`, `md`, `lg`)

### Card

A container component for displaying content in a contained format.

```jsx
import { Card, CardHeader, CardContent, CardFooter } from '@halooid/ui-components';

<Card elevation="md" hoverable>
  <CardHeader>
    <h3 className="card-title">Card Title</h3>
    <p className="card-subtitle">Card Subtitle</p>
  </CardHeader>
  <CardContent>
    <p>This is the card content.</p>
  </CardContent>
  <CardFooter>
    <Button variant="primary">Action</Button>
  </CardFooter>
</Card>
```

#### Props

- `elevation`: Card elevation (`none`, `sm`, `md`, `lg`, `xl`, `2xl`)
- `bordered`: Whether the card has a border
- `rounded`: Border radius (`none`, `sm`, `md`, `lg`, `xl`, `2xl`, `full`)
- `hoverable`: Whether the card is hoverable
- `onTap`: Callback when card is tapped

## Theme

The component library includes a theme system that provides consistent colors, typography, and spacing across all components.

```jsx
import { ThemeProvider } from 'styled-components';
import { theme } from '@halooid/ui-components';

// Use the default theme
<ThemeProvider theme={theme}>
  <App />
</ThemeProvider>

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

<ThemeProvider theme={customTheme}>
  <App />
</ThemeProvider>
```

### Theme Properties

- `colors`: Color palette
- `typography`: Typography settings
- `spacing`: Spacing function
- `breakpoints`: Responsive breakpoints
- `shadows`: Box shadow definitions
- `borderRadius`: Border radius definitions
- `transitions`: Animation transitions
- `zIndex`: Z-index values

## Utilities

### Hooks

- `useMediaQuery`: Hook for responsive design
- `useOutsideClick`: Hook for detecting clicks outside an element

### Functions

- `classNames`: Utility for combining class names

## Development

### Adding a New Component

1. Create a new directory in `src/components` with the component name
2. Create the component file (e.g., `ComponentName.tsx`)
3. Create a stories file (e.g., `ComponentName.stories.tsx`)
4. Export the component in `src/components/index.ts`
5. Document the component in Storybook

### Running Tests

```bash
# Run tests
npm test

# Run tests with coverage
npm test -- --coverage
```

### Building the Library

```bash
# Build the library
npm run build

# Watch for changes and rebuild
npm run dev
```

## Best Practices

1. **Consistency**: Follow the established patterns and styles
2. **Accessibility**: Ensure components are accessible
3. **Performance**: Keep components lightweight and efficient
4. **Documentation**: Document all props and usage examples
5. **Testing**: Write tests for all components
