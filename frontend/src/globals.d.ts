export {};

import * as runtime from '../wailsjs/runtime/runtime'

declare global {

    interface Window {
        runtime: typeof runtime;
        SendToast: (
            message: string,
            time: number,
            fade: number,
            bg?: string,
            col?: string,

        ) => void
    }
}