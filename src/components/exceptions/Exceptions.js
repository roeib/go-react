import React from 'react';
import ex1 from '../../assets/ex1.png'
import ex2 from '../../assets/ex2.png'
import ex3 from '../../assets/ex3.png'

const exceptionsImg = {
  DivideByZeroException: ex1,
  IOException: ex2,
  NullPointerException: ex3,
}
const Exceptions = ({ exceptions }) => {
    return (
        exceptions.map(exception => {
            return (
              <div key={Math.random()} style={{ position: 'absolute', bottom: exception.y + 'px', left: exception.x + 'px' }}>
                <img src={exceptionsImg[exception.exceptionType]} alt="" />
              </div>
            )
          })
        )

}
export default React.memo(Exceptions)