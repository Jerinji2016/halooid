import React from 'react';
import { Meta, StoryObj } from '@storybook/react';
import Input from './Input';

const meta: Meta<typeof Input> = {
  title: 'Components/Input',
  component: Input,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['outlined', 'filled', 'standard'],
      description: 'The variant of the input',
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
      description: 'The size of the input',
    },
    disabled: {
      control: 'boolean',
      description: 'Whether the input is disabled',
    },
    fullWidth: {
      control: 'boolean',
      description: 'Whether the input is full width',
    },
    label: {
      control: 'text',
      description: 'The label of the input',
    },
    placeholder: {
      control: 'text',
      description: 'The placeholder of the input',
    },
    helperText: {
      control: 'text',
      description: 'Helper text to display below the input',
    },
    error: {
      control: 'text',
      description: 'Error message to display below the input',
    },
    onChange: { action: 'changed' },
  },
};

export default meta;
type Story = StoryObj<typeof Input>;

export const Default: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
  },
};

export const Outlined: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    variant: 'outlined',
  },
};

export const Filled: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    variant: 'filled',
  },
};

export const Standard: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    variant: 'standard',
  },
};

export const WithHelperText: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    helperText: 'We will never share your email with anyone else.',
  },
};

export const WithError: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    error: 'Please enter a valid email address.',
  },
};

export const Disabled: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    disabled: true,
  },
};

export const Small: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    size: 'sm',
  },
};

export const Medium: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    size: 'md',
  },
};

export const Large: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    size: 'lg',
  },
};

export const FullWidth: Story = {
  args: {
    label: 'Email',
    placeholder: 'Enter your email',
    fullWidth: true,
  },
};
