<script>
    import Blurrer from './components/Blurrer.svelte';
    import Main from './components/Main.svelte';
    import StardewStyleToast from './components/StardewStyleToast.svelte';
    import logBtn from './assets/images/log_btn.png';

    import { converters, setConverters } from './stores'
    import { objectToArray } from './util';
    import { ShowDebugLogsFolder } from '../wailsjs/go/main/App.js'



    window.runtime.EventsOnce('DOM_READY', (data) => {
        console.log(data)
    })

    window.runtime.EventsOnce('AVAILABLE_CONVERTERS', (data) => {
        console.log(data)
        console.log(objectToArray(data))
        setConverters(objectToArray(data))
    })

    window.runtime.EventsOn('UPDATE_CHECK_INFO', (data) => {
        const { state, message } = data.data

        if (state == 'none') return

        const color = state == 'available' ? '#9B0A16' : '#997E16'

        window.SendToast(
            message,
            state == 'available' ? 30000 : 15000,
            200,
            color,
            'white'
        )
    })

    let _sendToast = null;

    $: {
        if (window) window.SendToast = _sendToast
    }
</script>

<StardewStyleToast bind:SendToast={_sendToast}/>
<Blurrer />

<main>
    <Main />
</main>

<button 
    id="logs-btn"
    on:click={ShowDebugLogsFolder}>
    <img src={logBtn} alt="open logs" />
</button>

<style lang="scss">
    @use "data";

    main {
        z-index: 4;

        width: var(--main-width);
        height: var(--main-height);

        background-color: wheat;

        padding: 30px;

        @include data.border_image;
    }

    #logs-btn {
        img {
            width: 1.5em;
            height: 1.5em;
        }
        z-index: 2;
        padding: 0.25em;

        background-color: wheat;
        @include data.border_image(8px);

        position: absolute;
        bottom: 8px;
        left: 8px;

        &:hover {
            filter: brightness(0.8);
        }
    }

</style>
