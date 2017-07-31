//module TileTypes {
    export enum TileStates {
        // default gray tile
        Default = 'default',
        // success green tile
        Success = 'success',
        // indeterminate cyan tile
        Indeterminate = 'indeterminate',
        // warning orange tile
        Warning = 'warning',
        // default red tile
        Error = 'error'
    }

    export enum TileSize {
        // tile that takes 1x1 grid cell
        x1 = '1x',
        // tile that takes 2x1 grid cells (horizontally)
        x2 = '2x',
        // tile that takes 2x2 grid cells
        x4 = '4x'
    }

    export enum TileTypes {
        // Tile with one line of text
        Text = 'text',
        // Tile with header and a line of text
        TestStatus = 'text-status',
        // Tile with small header and a large line of text
        TextStatus2 = 'text-status-2',
        // Tile with header, a large line of text and a progress bar
        TextStatusBar = 'text-status-bar'
    }
//}
