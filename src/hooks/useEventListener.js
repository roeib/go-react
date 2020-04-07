import {useRef, useEffect } from 'react';

export const  useEventListener=(eventName, handler, element = window)=>{
    // Create a ref that stores handler
    const savedHandler = useRef();
    
    useEffect(() => {
      savedHandler.current = handler;
    }, [handler]);
  
    useEffect( () => {
        const isSupported = element && element.addEventListener;
        if (!isSupported) return;

        const eventListener = event => {savedHandler.current(event)}

        element.addEventListener(eventName, eventListener)
        
        return () => {
          element.removeEventListener(eventName, eventListener);
        };
      },
      [eventName, element] 
    );
  };