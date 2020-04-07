
import React, { useCallback } from 'react';
import './App.css';

import { useWebSocket } from './hooks/useWS'
import { useBodyBounderies } from './hooks/useBodyBounderies'
import  {ws, playerMoves} from './Utills.js/Utills'
import { useEventListener } from './hooks/useEventListener'
import Monkeys from './components/Monkeys'
function App() {
  const bodyBounderies = useBodyBounderies()
  const [players, sendMSG] = useWebSocket(ws, bodyBounderies)
  const showKeyCode = useCallback(({ key }) => {
      sendMSG(playerMoves[key])
    },[sendMSG]);

  useEventListener('keydown', showKeyCode);
  return (
    <>
    <Monkeys players={players} />
    </>
  );
}
export default App;