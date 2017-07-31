import * as React from 'react'
import * as cn from 'classnames'

import { ITile } from 'contracts'
import { TileSize, TileStates } from './enums'

export interface IProps extends ITile {

}

export default function Tile(props: IProps) {
    let tileSize, tileState;
    switch (props.size) {
        case TileSize.x2:
            tileSize = 'tile--size-x2';
            break;
        case TileSize.x4:
            tileSize = 'tile--size-x4';
            break;
        default:
            tileSize = null;
            break;
    }

    switch (props.state) {
        case TileStates.Success:
            tileState = 'tile--state-success';
            break;
        case TileStates.Indeterminate:
            tileState = 'tile--state-indeterminate';
            break;
        case TileStates.Warning:
            tileState = 'tile--state-warning';
            break;
        case TileStates.Error:
            tileState = 'tile--state-error';
            break;
        default:
            tileState = 'tile--state-default';
            break;
    }

    const now = Date.now();
    const updated = Date.parse(props.updated);
    const difference = new Date(now - updated);

    return (
        <div className={cn('tile__container', tileSize)}>
            <div className={cn('tile__element', tileState)}>
                <div className='tile__title'>{props.titleText}</div>
                <div className="tile__content">{props.descrText}</div>
                <div className="tile__footer">{`updated ${difference.getMinutes()} mins`}</div>
            </div>
        </div>
    );
}