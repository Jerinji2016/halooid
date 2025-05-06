import React, { HTMLAttributes, ReactNode } from 'react';
import styled, { css } from 'styled-components';
import { classNames } from '../../utils';

export interface CardProps extends HTMLAttributes<HTMLDivElement> {
  /**
   * Card content
   */
  children: ReactNode;
  /**
   * Whether the card has a shadow
   */
  elevation?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  /**
   * Whether the card has a border
   */
  bordered?: boolean;
  /**
   * Whether the card has rounded corners
   */
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
  /**
   * Whether the card is hoverable
   */
  hoverable?: boolean;
  /**
   * Additional CSS class names
   */
  className?: string;
}

const StyledCard = styled.div<Omit<CardProps, 'children' | 'className'>>`
  display: flex;
  flex-direction: column;
  position: relative;
  background-color: ${({ theme }) => theme.colors.background.paper};
  overflow: hidden;
  
  /* Elevation styles */
  ${({ theme, elevation = 'md' }) => {
    switch (elevation) {
      case 'none':
        return css`
          box-shadow: ${theme.shadows.none};
        `;
      case 'sm':
        return css`
          box-shadow: ${theme.shadows.sm};
        `;
      case 'lg':
        return css`
          box-shadow: ${theme.shadows.lg};
        `;
      case 'xl':
        return css`
          box-shadow: ${theme.shadows.xl};
        `;
      case '2xl':
        return css`
          box-shadow: ${theme.shadows['2xl']};
        `;
      case 'md':
      default:
        return css`
          box-shadow: ${theme.shadows.md};
        `;
    }
  }}
  
  /* Border styles */
  ${({ bordered, theme }) => bordered && css`
    border: 1px solid ${theme.colors.neutral.light};
  `}
  
  /* Rounded styles */
  ${({ theme, rounded = 'md' }) => {
    switch (rounded) {
      case 'none':
        return css`
          border-radius: ${theme.borderRadius.none};
        `;
      case 'sm':
        return css`
          border-radius: ${theme.borderRadius.sm};
        `;
      case 'lg':
        return css`
          border-radius: ${theme.borderRadius.lg};
        `;
      case 'xl':
        return css`
          border-radius: ${theme.borderRadius.xl};
        `;
      case '2xl':
        return css`
          border-radius: ${theme.borderRadius['2xl']};
        `;
      case 'full':
        return css`
          border-radius: ${theme.borderRadius.full};
        `;
      case 'md':
      default:
        return css`
          border-radius: ${theme.borderRadius.md};
        `;
    }
  }}
  
  /* Hoverable styles */
  ${({ hoverable, theme }) => hoverable && css`
    transition: box-shadow 0.3s ease-in-out, transform 0.3s ease-in-out;
    
    &:hover {
      box-shadow: ${theme.shadows.lg};
      transform: translateY(-2px);
    }
  `}
`;

export const CardHeader = styled.div`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid ${({ theme }) => theme.colors.neutral.light};
  
  .card-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: ${({ theme }) => theme.colors.text.primary};
  }
  
  .card-subtitle {
    margin: 0.25rem 0 0;
    font-size: 0.875rem;
    color: ${({ theme }) => theme.colors.text.secondary};
  }
`;

export const CardContent = styled.div`
  padding: 1rem;
  flex: 1 1 auto;
`;

export const CardFooter = styled.div`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-top: 1px solid ${({ theme }) => theme.colors.neutral.light};
`;

export const CardMedia = styled.div<{ image?: string; height?: string }>`
  background-image: ${({ image }) => image ? `url(${image})` : 'none'};
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  height: ${({ height }) => height || '200px'};
`;

/**
 * Card component for displaying content in a contained format
 */
export const Card: React.FC<CardProps> = ({
  children,
  elevation = 'md',
  bordered = false,
  rounded = 'md',
  hoverable = false,
  className,
  ...props
}) => {
  return (
    <StyledCard
      elevation={elevation}
      bordered={bordered}
      rounded={rounded}
      hoverable={hoverable}
      className={classNames('halooid-card', className || '')}
      {...props}
    >
      {children}
    </StyledCard>
  );
};

export default Card;
