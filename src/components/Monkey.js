import React from 'react';
import monkey from '../assets/monkey.png'
import activemonkey from '../assets/activemonkey.png'
const Monkey = ({ player }) => {
    return (
        <>
            <div className={`playerImg ${player.shake === true ? "shake" : ''}`} style={{ bottom: player.p.y + 'px', left: player.p.x + 'px', color: `rgb(${player.color[0]},${player.color[1]},${player.color[2]})` }}>
                <span style={{ position: "absolute", top: "-40px", left: "-25px" }}>{player.exceptionType}score{player.score}</span>
                <img alt="player img" src={player.active ? activemonkey : monkey} />
            </div>
        </>

    );
}
export default Monkey