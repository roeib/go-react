
import React, { useCallback } from 'react';
import './App.css';
import { useWebSocket } from './hooks/useWebSocket'
import { useBodyBounderies } from './hooks/useBodyBounderies'
import { ws, playerMoves } from './Utils.js/Utils'
import { useEventListener } from './hooks/useEventListener'
import Monkeys from './components/Monkeys'
import ex1 from './assets/ex1.png'
import ex2 from './assets/ex2.png'
import ex3 from './assets/ex3.png'

const exceptionsImg = {
  DivideByZeroException: ex1,
  IOException: ex2,
  NullPointerException: ex3,
}
function App() {
  const bodyBounderies = useBodyBounderies()
  const [playState, sendMSG] = useWebSocket(ws, bodyBounderies)
  const showKeyCode = useCallback(({ key }) => {
    sendMSG(playerMoves[key])
  }, [sendMSG]);

  useEventListener('keydown', showKeyCode);
  return (
    <>
      <div>
        <h2>Score</h2>
        {playState.players.map(player => <div key={player.id}>{player.active ? "You" : player.exceptionType} - {player.score}</div>)}
      </div>
      {
        playState.exceptions.map(exception => {
          return (
            <div key={Math.random()} style={{ position: 'absolute', bottom: exception.p.y + 'px', left: exception.p.x + 'px' }}>
              <img src={exceptionsImg[exception.exceptionType]} alt=""/>
            </div>
          )
        })
      }

      <Monkeys players={playState.players} />
    </>
  );
}
export default App;