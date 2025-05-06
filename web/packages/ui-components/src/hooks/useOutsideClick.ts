import { useEffect, RefObject } from 'react';

/**
 * Hook that triggers a callback when a click occurs outside of the specified element
 * 
 * @param ref - Reference to the element to detect clicks outside of
 * @param callback - Function to call when a click outside occurs
 * 
 * @example
 * const ref = useRef(null);
 * useOutsideClick(ref, () => setIsOpen(false));
 */
export function useOutsideClick<T extends HTMLElement = HTMLElement>(
  ref: RefObject<T>,
  callback: () => void
): void {
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        callback();
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [ref, callback]);
}
