import React from 'react';
import monkey from '../assets/monkey.png'
import activemonkey from '../assets/activemonkey.png'
const Monkey = ({ shake,p,color,exceptionType,active,score }) => {
    return (
        <>
            <div className={`playerImg ${shake === true ? "shake" : ''}`} style={{ bottom: p.y + 'px', left: p.x + 'px', color: `rgb(${color[0]},${color[1]},${color[2]})` }}>
                <span style={{ position: "absolute", top: "-40px", left: "-25px" }}>{exceptionType}score{score}</span>
                <img alt="player img" src={active ? activemonkey : monkey} />
            </div>
        </>

    );
}
export default Monkey