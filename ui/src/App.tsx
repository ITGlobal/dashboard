import * as React from 'react'
import * as cn from 'classnames'
import * as _ from 'lodash'

import { IDataJson, ITile } from 'contracts'
import { TileStates, TileTypes } from './enums'
import { getData } from './api'
import TileComponent from './Tile'

import * as data from './data.json'

import './style/style.scss';

export interface IProps {

}

interface IAppState {
    items: ITile[];
}

class App extends React.Component<IProps, IAppState> {
    constructor(props) {
        super(props)

        this.state = {
            items: []
        }

        setInterval(() => {
            this.getTiles();
        }, 5000)

        this.getStateNum = this.getStateNum.bind(this);
    }

    getTiles() {
        // let i;
        // getData<IDataJson>().then(resp => i = resp.tiles);
        // console.log(i);
        this.setState({ items: data['tiles'] });
    }

    getStateNum(tile: ITile): number {
        let res = 0;
        switch (tile.state) {
            case TileStates.Error:
                res += 0;
                break;
            case TileStates.Warning:
                res += 100;
                break;
            case TileStates.Indeterminate:
                res += 200;
                break;
            case TileStates.Success:
                res += 300;
                break;
            default: res += 900;
                break;
        }

        switch (tile.type) {
            case TileTypes.Text:
                res += 0;
                break;
            case TileTypes.TextStatus:
                res += 2000;
                break;
            case TileTypes.TextStatus2:
                res += 1000;
                break;
            case TileTypes.TextStatusBar:
                res += 3000;
                break;
            default: 
                res += 9000;
                break;
        }

        return res;
    }

    componentDidMount() {
        this.getTiles();
    }

    render() {
        const { items } = this.state;
        const groups = Object.values(_.groupBy(items, 'type'));
        return (
            <div className="tiles__container">
                {
                    items.length > 0 ?
                        items
                            .sort((x: ITile, y: ITile) => x.titleText.localeCompare(y.titleText))
                            .sort((x: ITile, y: ITile) => this.getStateNum(x) - this.getStateNum(y))
                            .map((item) => <TileComponent {...item} key={item.id} />)
                        : ''
                }
            </div>
        )
    };
}

export default App
