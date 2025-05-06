import React from 'react';
import { Meta, StoryObj } from '@storybook/react';
import Card, { CardHeader, CardContent, CardFooter, CardMedia } from './Card';
import Button from '../Button';

const meta: Meta<typeof Card> = {
  title: 'Components/Card',
  component: Card,
  tags: ['autodocs'],
  argTypes: {
    elevation: {
      control: 'select',
      options: ['none', 'sm', 'md', 'lg', 'xl', '2xl'],
      description: 'The elevation of the card',
    },
    bordered: {
      control: 'boolean',
      description: 'Whether the card has a border',
    },
    rounded: {
      control: 'select',
      options: ['none', 'sm', 'md', 'lg', 'xl', '2xl', 'full'],
      description: 'The border radius of the card',
    },
    hoverable: {
      control: 'boolean',
      description: 'Whether the card is hoverable',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Card>;

export const Default: Story = {
  args: {
    children: (
      <CardContent>
        <p>This is a basic card with content.</p>
      </CardContent>
    ),
  },
};

export const WithHeaderAndFooter: Story = {
  args: {
    children: (
      <>
        <CardHeader>
          <div>
            <h3 className="card-title">Card Title</h3>
            <p className="card-subtitle">Card Subtitle</p>
          </div>
        </CardHeader>
        <CardContent>
          <p>This is a card with a header and footer.</p>
          <p>You can put any content here.</p>
        </CardContent>
        <CardFooter>
          <Button variant="primary" size="sm">Action</Button>
          <Button variant="secondary" size="sm" outlined style={{ marginLeft: '0.5rem' }}>Cancel</Button>
        </CardFooter>
      </>
    ),
  },
};

export const WithMedia: Story = {
  args: {
    children: (
      <>
        <CardMedia 
          image="https://images.unsplash.com/photo-1522252234503-e356532cafd5?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=800&q=80" 
          height="200px"
        />
        <CardHeader>
          <div>
            <h3 className="card-title">Card with Media</h3>
            <p className="card-subtitle">A card with an image</p>
          </div>
        </CardHeader>
        <CardContent>
          <p>This card has a media section at the top.</p>
        </CardContent>
      </>
    ),
  },
};

export const Bordered: Story = {
  args: {
    bordered: true,
    elevation: 'none',
    children: (
      <CardContent>
        <p>This is a bordered card without elevation.</p>
      </CardContent>
    ),
  },
};

export const Hoverable: Story = {
  args: {
    hoverable: true,
    children: (
      <CardContent>
        <p>Hover over this card to see the effect.</p>
      </CardContent>
    ),
  },
};

export const RoundedCorners: Story = {
  args: {
    rounded: 'lg',
    children: (
      <CardContent>
        <p>This card has large rounded corners.</p>
      </CardContent>
    ),
  },
};
