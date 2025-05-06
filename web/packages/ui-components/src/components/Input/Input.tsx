import React, { InputHTMLAttributes, ReactNode, forwardRef } from 'react';
import styled, { css } from 'styled-components';
import { classNames } from '../../utils';

export type InputSize = 'sm' | 'md' | 'lg';
export type InputVariant = 'outlined' | 'filled' | 'standard';

export interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  /**
   * Input label
   */
  label?: string;
  /**
   * Input size
   */
  size?: InputSize;
  /**
   * Input variant
   */
  variant?: InputVariant;
  /**
   * Error message
   */
  error?: string;
  /**
   * Helper text
   */
  helperText?: string;
  /**
   * Whether the input is full width
   */
  fullWidth?: boolean;
  /**
   * Icon to display at the start of the input
   */
  startIcon?: ReactNode;
  /**
   * Icon to display at the end of the input
   */
  endIcon?: ReactNode;
  /**
   * Additional CSS class names
   */
  className?: string;
}

const InputContainer = styled.div<{
  size?: InputSize;
  variant?: InputVariant;
  hasError?: boolean;
  disabled?: boolean;
  fullWidth?: boolean;
}>`
  display: inline-flex;
  flex-direction: column;
  position: relative;
  font-family: ${({ theme }) => theme.typography.fontFamily};
  
  ${({ fullWidth }) => fullWidth && css`
    width: 100%;
  `}
  
  .input-label {
    margin-bottom: 0.25rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: ${({ theme, hasError }) => hasError ? theme.colors.error.main : theme.colors.text.primary};
  }
  
  .input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    
    ${({ size }) => {
      switch (size) {
        case 'sm':
          return css`
            height: 2rem;
            font-size: 0.75rem;
          `;
        case 'lg':
          return css`
            height: 3rem;
            font-size: 1rem;
          `;
        case 'md':
        default:
          return css`
            height: 2.5rem;
            font-size: 0.875rem;
          `;
      }
    }}
    
    .start-icon, .end-icon {
      position: absolute;
      display: flex;
      align-items: center;
      justify-content: center;
      color: ${({ theme }) => theme.colors.text.secondary};
      z-index: 1;
    }
    
    .start-icon {
      left: 0.75rem;
    }
    
    .end-icon {
      right: 0.75rem;
    }
  }
  
  .input-field {
    width: 100%;
    height: 100%;
    padding: 0 0.75rem;
    border: none;
    outline: none;
    background-color: transparent;
    font-family: inherit;
    font-size: inherit;
    color: ${({ theme }) => theme.colors.text.primary};
    
    &::placeholder {
      color: ${({ theme }) => theme.colors.text.secondary};
    }
    
    ${({ startIcon }) => startIcon && css`
      padding-left: 2.25rem;
    `}
    
    ${({ endIcon }) => endIcon && css`
      padding-right: 2.25rem;
    `}
  }
  
  /* Variant styles */
  ${({ theme, variant, hasError }) => {
    switch (variant) {
      case 'filled':
        return css`
          .input-wrapper {
            background-color: ${theme.colors.neutral.light};
            border-radius: ${theme.borderRadius.md};
            
            &:hover {
              background-color: ${theme.colors.neutral.light};
            }
            
            &:focus-within {
              background-color: ${theme.colors.background.paper};
              box-shadow: 0 0 0 2px ${hasError ? theme.colors.error.main : theme.colors.primary.main};
            }
          }
        `;
      case 'standard':
        return css`
          .input-wrapper {
            border-bottom: 1px solid ${hasError ? theme.colors.error.main : theme.colors.neutral.main};
            
            &:hover {
              border-bottom: 1px solid ${hasError ? theme.colors.error.main : theme.colors.text.primary};
            }
            
            &:focus-within {
              border-bottom: 2px solid ${hasError ? theme.colors.error.main : theme.colors.primary.main};
            }
          }
          
          .input-field {
            padding-left: 0;
            padding-right: 0;
          }
          
          ${({ startIcon }) => startIcon && css`
            .input-field {
              padding-left: 1.5rem;
            }
            
            .start-icon {
              left: 0;
            }
          `}
          
          ${({ endIcon }) => endIcon && css`
            .input-field {
              padding-right: 1.5rem;
            }
            
            .end-icon {
              right: 0;
            }
          `}
        `;
      case 'outlined':
      default:
        return css`
          .input-wrapper {
            border: 1px solid ${hasError ? theme.colors.error.main : theme.colors.neutral.main};
            border-radius: ${theme.borderRadius.md};
            
            &:hover {
              border-color: ${hasError ? theme.colors.error.main : theme.colors.text.primary};
            }
            
            &:focus-within {
              border-color: ${hasError ? theme.colors.error.main : theme.colors.primary.main};
              box-shadow: 0 0 0 1px ${hasError ? theme.colors.error.main : theme.colors.primary.main};
            }
          }
        `;
    }
  }}
  
  /* Disabled style */
  ${({ disabled, theme }) => disabled && css`
    opacity: 0.7;
    cursor: not-allowed;
    
    .input-label {
      color: ${theme.colors.text.disabled};
    }
    
    .input-wrapper {
      background-color: ${theme.colors.neutral.light};
      border-color: ${theme.colors.neutral.main};
      
      &:hover {
        border-color: ${theme.colors.neutral.main};
      }
    }
    
    .input-field {
      cursor: not-allowed;
      color: ${theme.colors.text.disabled};
      
      &::placeholder {
        color: ${theme.colors.text.disabled};
      }
    }
    
    .start-icon, .end-icon {
      color: ${theme.colors.text.disabled};
    }
  `}
  
  .helper-text {
    margin-top: 0.25rem;
    font-size: 0.75rem;
    color: ${({ theme, hasError }) => hasError ? theme.colors.error.main : theme.colors.text.secondary};
  }
`;

/**
 * Input component for user text input
 */
export const Input = forwardRef<HTMLInputElement, InputProps>(({
  label,
  size = 'md',
  variant = 'outlined',
  error,
  helperText,
  fullWidth = false,
  startIcon,
  endIcon,
  className,
  disabled,
  ...props
}, ref) => {
  const hasError = !!error;
  const helpText = error || helperText;
  
  return (
    <InputContainer
      size={size}
      variant={variant}
      hasError={hasError}
      disabled={disabled}
      fullWidth={fullWidth}
      startIcon={!!startIcon}
      endIcon={!!endIcon}
      className={classNames('halooid-input', className || '')}
    >
      {label && (
        <label className="input-label">
          {label}
        </label>
      )}
      <div className="input-wrapper">
        {startIcon && (
          <span className="start-icon">
            {startIcon}
          </span>
        )}
        <input
          ref={ref}
          className="input-field"
          disabled={disabled}
          {...props}
        />
        {endIcon && (
          <span className="end-icon">
            {endIcon}
          </span>
        )}
      </div>
      {helpText && (
        <div className="helper-text">
          {helpText}
        </div>
      )}
    </InputContainer>
  );
});

Input.displayName = 'Input';

export default Input;
