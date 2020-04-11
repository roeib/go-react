
import React, { useCallback } from 'react';
import './App.css';
import { useWebSocket } from './hooks/useWebSocket'
import { useBodyBounderies } from './hooks/useBodyBounderies'
import {ws, playerMoves} from './Utills.js/Utills'
import { useEventListener } from './hooks/useEventListener'
import Monkeys from './components/Monkeys'

function App() {
  const bodyBounderies = useBodyBounderies()
  const [playState, sendMSG] = useWebSocket(ws, bodyBounderies)
  const showKeyCode = useCallback(({ key }) => {
      sendMSG(playerMoves[key])
    },[sendMSG]);

  useEventListener('keydown', showKeyCode);
  return (
    <>
    <div>
      <h2>Score</h2>
      {playState.players.map(player => <div>{player.id} - {player.score}</div>)}
    </div>
    <Monkeys players={playState.players} />
    </>
  );
}
export default App;