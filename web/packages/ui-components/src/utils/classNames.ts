/**
 * Combines multiple class names into a single string
 * Filters out falsy values
 * 
 * @param classes - Class names to combine
 * @returns Combined class names as a string
 * 
 * @example
 * classNames('btn', 'btn-primary', condition && 'active')
 * // Returns: 'btn btn-primary active' if condition is true
 * // Returns: 'btn btn-primary' if condition is false
 */
export function classNames(...classes: (string | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(' ');
}
