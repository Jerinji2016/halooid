import React, { ButtonHTMLAttributes, ReactNode } from 'react';
import styled, { css } from 'styled-components';
import { classNames } from '../../utils';

export type ButtonVariant = 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info' | 'light' | 'dark' | 'link';
export type ButtonSize = 'sm' | 'md' | 'lg';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /**
   * Button variant
   */
  variant?: ButtonVariant;
  /**
   * Button size
   */
  size?: ButtonSize;
  /**
   * Whether the button is outlined
   */
  outlined?: boolean;
  /**
   * Whether the button is rounded
   */
  rounded?: boolean;
  /**
   * Whether the button is disabled
   */
  disabled?: boolean;
  /**
   * Whether the button is full width
   */
  fullWidth?: boolean;
  /**
   * Button content
   */
  children: ReactNode;
  /**
   * Icon to display before the button text
   */
  startIcon?: ReactNode;
  /**
   * Icon to display after the button text
   */
  endIcon?: ReactNode;
  /**
   * Whether the button is in a loading state
   */
  loading?: boolean;
  /**
   * Additional CSS class names
   */
  className?: string;
}

const StyledButton = styled.button<ButtonProps>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  box-sizing: border-box;
  outline: 0;
  border: 0;
  cursor: pointer;
  user-select: none;
  vertical-align: middle;
  text-decoration: none;
  font-weight: ${({ theme }) => theme.typography.fontWeightMedium};
  font-family: ${({ theme }) => theme.typography.fontFamily};
  border-radius: ${({ theme }) => theme.borderRadius.md};
  transition: background-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms,
    box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms,
    border-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms,
    color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms;
  
  /* Size styles */
  ${({ size }) => {
    switch (size) {
      case 'sm':
        return css`
          padding: 0.25rem 0.5rem;
          font-size: 0.75rem;
        `;
      case 'lg':
        return css`
          padding: 0.75rem 1.5rem;
          font-size: 1.125rem;
        `;
      case 'md':
      default:
        return css`
          padding: 0.5rem 1rem;
          font-size: 0.875rem;
        `;
    }
  }}
  
  /* Variant styles */
  ${({ theme, variant, outlined }) => {
    const variantColor = variant && theme.colors[variant as keyof typeof theme.colors]
      ? theme.colors[variant as keyof typeof theme.colors]
      : theme.colors.primary;
    
    if (outlined) {
      return css`
        color: ${variantColor.main};
        border: 1px solid ${variantColor.main};
        background-color: transparent;
        
        &:hover {
          background-color: rgba(${variantColor.main}, 0.04);
        }
        
        &:active {
          background-color: rgba(${variantColor.main}, 0.12);
        }
      `;
    }
    
    return css`
      color: ${variantColor.contrastText};
      background-color: ${variantColor.main};
      
      &:hover {
        background-color: ${variantColor.dark};
      }
      
      &:active {
        background-color: ${variantColor.dark};
      }
    `;
  }}
  
  /* Rounded style */
  ${({ rounded, theme }) => rounded && css`
    border-radius: ${theme.borderRadius.full};
  `}
  
  /* Full width style */
  ${({ fullWidth }) => fullWidth && css`
    width: 100%;
  `}
  
  /* Disabled style */
  ${({ disabled, theme }) => disabled && css`
    color: ${theme.colors.text.disabled};
    background-color: ${theme.colors.neutral.light};
    cursor: default;
    pointer-events: none;
    
    &:hover, &:active {
      background-color: ${theme.colors.neutral.light};
    }
  `}
  
  /* Loading style */
  ${({ loading }) => loading && css`
    cursor: wait;
    
    .button-text {
      opacity: 0;
    }
    
    .loading-indicator {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
    }
  `}
  
  /* Icon spacing */
  .start-icon {
    margin-right: 0.5rem;
  }
  
  .end-icon {
    margin-left: 0.5rem;
  }
`;

/**
 * Button component for user interaction
 */
export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  outlined = false,
  rounded = false,
  disabled = false,
  fullWidth = false,
  children,
  startIcon,
  endIcon,
  loading = false,
  className,
  ...props
}) => {
  return (
    <StyledButton
      variant={variant}
      size={size}
      outlined={outlined}
      rounded={rounded}
      disabled={disabled || loading}
      fullWidth={fullWidth}
      className={classNames('halooid-button', className || '')}
      {...props}
    >
      {startIcon && <span className="start-icon">{startIcon}</span>}
      <span className="button-text">{children}</span>
      {endIcon && <span className="end-icon">{endIcon}</span>}
      {loading && (
        <span className="loading-indicator">
          {/* Simple loading spinner */}
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <style>
              {`
                .spinner {
                  animation: rotate 2s linear infinite;
                  transform-origin: center;
                }
                .path {
                  stroke-linecap: round;
                  animation: dash 1.5s ease-in-out infinite;
                }
                @keyframes rotate {
                  100% {
                    transform: rotate(360deg);
                  }
                }
                @keyframes dash {
                  0% {
                    stroke-dasharray: 1, 150;
                    stroke-dashoffset: 0;
                  }
                  50% {
                    stroke-dasharray: 90, 150;
                    stroke-dashoffset: -35;
                  }
                  100% {
                    stroke-dasharray: 90, 150;
                    stroke-dashoffset: -124;
                  }
                }
              `}
            </style>
            <circle
              className="spinner"
              cx="12"
              cy="12"
              r="10"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            />
          </svg>
        </span>
      )}
    </StyledButton>
  );
};

export default Button;
