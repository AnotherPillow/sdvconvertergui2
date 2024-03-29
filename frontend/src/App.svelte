<script>
    import Blurrer from './components/Blurrer.svelte';
    import Main from './components/Main.svelte';
    import StardewStyleToast from './components/StardewStyleToast.svelte';

    import { converters, setConverters } from './stores'
    import { objectToArray } from './util';


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

</style>
