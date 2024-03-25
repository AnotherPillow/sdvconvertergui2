export const objectToArray = (dict) => {
    let output = []

    for (const key of Object.keys(dict)) {
        output.push(dict[key])
    }

    return output
}


export const catchAndNotify = (fun, message) => {
    try {
        return fun()
    }  catch {
        window.SendToast?.(message, 2000, 100, '#f5bbb3')
    }
}


// BROKEN - TODO: Fix
export const downloadLocalFile = (filepath) => {
    const url = `file:///${filepath.replaceAll('\\', '/')}`

    const tab = open(url)
    
    const timeout = (n = 5000) => setTimeout(() => {
        if (tab.window)
            return tab.window.close()
        else 
            return timeout(2000)
    }, n)
    timeout()
    
}