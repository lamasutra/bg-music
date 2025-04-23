<template>
  <q-page>
    <q-header reveal elevated>
      <music-player title="Music player" />
    </q-header>
    <div v-if="showOptions" style="padding-top: 44px;">
      <options-panel/>
      <q-page-sticky expand position="top" :offset="[0, 0]">
        <q-toolbar class="bg-dark">
          <q-toolbar-title>Options</q-toolbar-title>
        </q-toolbar>
      </q-page-sticky>
    </div>
  </q-page>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { EventsOn } from '../../wailsjs/runtime/runtime'
  import { SetWindowSize, GetWindowSize } from '../../wailsjs/go/ui/guiState'
  import MusicPlayer from 'components/MusicPlayer.vue'
  import OptionsPanel from 'components/OptionsPanel.vue'
  import { ui } from '../../wailsjs/go/models'

  const showOptions = ref(false)
  onMounted(() => {
    var defaultSize: ui.WindowSize
    GetWindowSize().then(
      (size: ui.WindowSize) => {
        defaultSize = size
        console.info(size)
      }
    )

    EventsOn('show-options', (show: boolean) => {
      showOptions.value = show
      if (show) {
        SetWindowSize(defaultSize.width, defaultSize.height+300)
      } else {
        SetWindowSize(defaultSize.width, defaultSize.height)
      }
    })
  })
</script>
