import { TileTypes, TileSize, TileStates } from 'enums'

/**
 * Интерфейс тайла
 */
export interface ITile {
    /** Tile ID. It's is stable and should not change */
    id: string;
    /** Corresponding tile provider type */
    source: string;
    /** Tile's last update time */
    updated: string;
    /** Tile layout type */
    type: TileTypes;
    /** Tile size */
    size: TileSize;
    /** Tile state */
    state: TileStates;
    /** Tile title text */
    titleText: string;
    /** Tile description text */
    descrText?: string;
    /** Tile's progress bar position */
    statusValue?: number;
}

export interface IDataJson {
    tiles: ITile[];
    version: string;
}
