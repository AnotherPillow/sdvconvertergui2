import { writable } from "svelte/store";

/**@type {import("svelte/store").Writable<{name: string; needs16: boolean; repo: string; downloadPath: string; inputDirectory: string; outputDirectory: string}[]>} */
export const converters = writable([])
/**
 * 
 * @param {any[]} array 
 */
export const setConverters = (array) => {
    if (array.length >= 1) converters.set(array)
}

export const dialogOpen = writable(false)