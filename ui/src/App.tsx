import * as React from 'react'
import * as cn from 'classnames'
import * as _ from 'lodash'

import { IDataJson, ITile } from 'contracts'
import { TileStates } from './enums'
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
        //getData<IDataJson>().then(resp => items = resp.tiles)
        this.setState({ items: data['tiles'] });
    }

    getStateNum(state: String) {
        switch (state) {
            case TileStates.Error: return 0;
            case TileStates.Warning: return 10;
            case TileStates.Indeterminate: return 20;
            case TileStates.Success: return 30;
            default: return 100;
        }
    }

    componentDidMount() {
        this.getTiles();
    }

    render() {
        const { items } = this.state;
        const groups = Object.values(_.groupBy(items, 'type'));
        return (
            <div className="dashboard__container">
                {
                    items.length > 0 ?
                        groups.map((group, i) => {
                            return (<div className="tiles__container" key={i}>{
                                group ?
                                    group
                                        .sort((x: ITile, y: ITile) => this.getStateNum(x.state) - this.getStateNum(y.state))
                                        .map((item) => <TileComponent {...item} key={item.id} />)
                                    : ''
                            }</div>)
                        })
                        :
                        ''
                }
            </div>
        )
    };
}

export default App