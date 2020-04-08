import React from 'react';
import monkey from '../assets/monkey.png'
import activemonkey from '../assets/activemonkey.png'
const Monkey = ({ shake,p,color,exceptionType,active,score }) => {
    return (
        <>
            <div className={`playerImg ${shake === true ? "shake" : ''}`} style={{ bottom: p.y + 'px', left: p.x + 'px' }}>
                <span className="exception">{exceptionType}</span>
                <img alt="player img" src={active ? activemonkey : monkey} />
            </div>
        </>

    );
}
export default Monkey