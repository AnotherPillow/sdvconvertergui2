<script lang="ts">

    import convertBtn from '../assets/images/convert_btn.png';
    import exportBtn from '../assets/images/export_btn.png';
    import ConverterOutput from './ConverterOutput.svelte';
    import ConvertersCombobox from './ConvertersCombobox.svelte';
    import UploadFiles from './UploadFiles.svelte';
    import { ConvertMod, ShowFolderInExplorer } from '../../wailsjs/go/main/App.js'
    import { downloadLocalFile } from '../util';

    let selectedConverter = '';
    // /** @type {{filename: string | null; content: string | null} | null}*/
    /** @type {any | null}*/
    let selectedManifest = null;
    let selectedFilePath = '';

    let conversionComplete = false;
    let downloadFileURL = ''

    window.runtime.EventsOn('CONVERTER_MOD_DONE', (path) => {
        conversionComplete = true;
        downloadFileURL = path
    })

    
    const exportMod = () => {
        window.SendToast('The zipped mod has been opened in your file manager.', 2000, 200, 'rgb(168, 190, 226);')
        ShowFolderInExplorer(downloadFileURL)
    }

    const convertMod = () => {
        conversionComplete = false;
        downloadFileURL = '';

        if (!selectedManifest) window.SendToast(
            'No manifest selected!', 2000, 100, '#f5bbb3'
        )

        if (!selectedConverter) window.SendToast(
            'No converter selected!', 2000, 100, '#f5bbb3'
        )

        ConvertMod(selectedManifest, selectedConverter, selectedFilePath).then(response => {
            const [type, message] = response.split('|')

            const colour = type == 'error' ? '#f5bbb3' : '#D1F5B3'

            window.SendToast(message, 6000, 200, colour)
        })
    }

</script>

<h1 class="title m-0">sdvconvertergui2</h1>
<div id="columns">
    <div class="column col-right">
        <!-- <h3 class="selected-mod-name"> -->
        <UploadFiles bind:selectedManifest={selectedManifest} bind:selectedFilePath={selectedFilePath} />
        <button id="convert-btn" on:click={convertMod}>
            <img src={convertBtn} alt="convert" id="convert-img" />
        </button>
        {#if conversionComplete}
            <button id="export-btn" on:click={exportMod}>
                <img src={exportBtn} alt="convert" id="export-img"/>
            </button>
        {/if}
    </div>
    <div class="column">
        <ConvertersCombobox bind:selectedValue={selectedConverter} />
        <ConverterOutput />
    </div>
</div>


<style lang="scss">
    @use "../data";
    h1 {
        font-weight: bold;
        font-size: 2em;
    }

    .result {
        height: 20px;
        line-height: 20px;
        margin: 1.5rem auto;
    }

    /*.input-box .btn {
        width: 60px;
        height: 30px;
        line-height: 30px;
        border-radius: 3px;
        border: none;
        margin: 0 0 0 20px;
        padding: 0 8px;
        cursor: pointer;
    }

    .input-box .btn:hover {
        background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
        color: #333333;
    }

    .input-box .input {
        border: none;
        border-radius: 3px;
        outline: none;
        height: 30px;
        line-height: 30px;
        padding: 0 10px;
        background-color: rgba(240, 240, 240, 1);
        -webkit-font-smoothing: antialiased;
    }

    .input-box .input:hover {
        border: none;
        background-color: rgba(255, 255, 255, 1);
    }

    .input-box .input:focus {
        border: none;
        background-color: rgba(255, 255, 255, 1);
    }*/

    #columns {
        display: flex;
        justify-content: space-between;
        align-items: center;

        height: data.$available-col-height;

        gap: 20px;

        .col-right {
            justify-content: space-evenly
        }

        .column {
            width: 50%;
            height: calc(100% - 60px);
            
            padding: 1em;

            display: flex;
            position: relative;
            flex-direction: column;
            
            align-items: center;

            @include data.border_image(12px, "../");
        }
    }

    #convert-btn, #export-btn {
        background-color: none;
        
        outline: none;
        border: none;

        &:hover {
            filter: brightness(0.8);
        }
        
        #convert-img, #export-img {
            height: 3em;
            aspect-ratio: 3 / 1;
        }
    }
</style>