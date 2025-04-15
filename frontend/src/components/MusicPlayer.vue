<template>
  <q-card class="q-mx-auto music-player no-shadow" flat square >
    <q-card-actions align="right" class="window-actions">
        <q-btn flat icon="minimize" @click="hide" title="Hide" />
        <q-btn flat icon="close" @click="close" title="Close"/>
      </q-card-actions>

    <q-card-section>
      <div class="title text-h6" title="track" style="--wails-draggable:drag">{{ title }}</div>
      <div class="text-subtitle2 text-grey" title="artist">{{ artist }}</div>
    </q-card-section>

    <q-card-section>
      <q-linear-progress
        :value="progress"
        color="primary"
        track-color="grey-4"
        size="10px"
        title="Progress"
        rounded
        instant-feedback
      />
      <div class="row justify-between q-mt-xs text-caption">
        <span title="Progress">{{ formatTime(currentTime) }}</span>
        <span title="Duration">{{ formatTime(duration) }}</span>
      </div>
    </q-card-section>

    <q-card-actions align="around">
      <q-btn flat round icon="fast_rewind" @click="rewind" title="Previous"/>
      <q-btn flat round :icon="isPlaying ? 'pause' : 'play_arrow'" @click="togglePlayback" :title="isPlaying ? 'Pause' : 'Play'" />
      <q-btn flat round icon="fast_forward" @click="forward" title="Next"/>
    </q-card-actions>

    <q-card-section>
      <div class="row items-center q-gutter-sm">
        <q-icon name="volume_up" />
        <q-slider
          v-model="volume"
          :min="0"
          :max="1"
          :step="0.01"
          color="primary"
          @change="updateVolume"
        />
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { SetVolume, Next, Prev, Pause, GetCurrentMetadata, GetCurrentMusic } from '../../wailsjs/go/gui/PlayerControls'
import { ToggleVisibility, CloseApp } from 'app/wailsjs/go/ui/guiState'
import { model } from '../../wailsjs/go/models'
// Make sure the Wails runtime is available
declare global {
  interface Window {
    runtime: any
  }
}

// Music meta
const title = ref<string>('')
const artist = ref<string>('')

// Audio setup
const isPlaying = ref<boolean>(true)
const currentTime = ref<number>(0)
const duration = ref<number>(0)
const progress = ref<number>(0)
const volume = ref<number>(1)

const hide = () => {
  // hides window
  ToggleVisibility()
}
const close = () => {
  // closes app
  CloseApp()
}

// Playback controls
const togglePlayback = () => {
  Pause()
}

const rewind = () => {
  Prev()
}

const forward = () => {
  Next()
}

const updateVolume = () => {
  console.info('set volume', volume.value)

  SetVolume(Math.round(volume.value * 100))
}

onMounted(() => {
  const evl1 = window.runtime.EventsOn('player_progress', (data: number) => {
    progress.value = data
  })
  const evl2 = window.runtime.EventsOn('player_music', (music: model.Music, metadata: model.MusicMetadata) => {
    duration.value = metadata.duration
    title.value = metadata.title || music.path || 'Unknown'
    artist.value = metadata.artist || 'Unknown'
  })
  const evl3 = window.runtime.EventsOn('player_volume', (vol: number) => {
    volume.value = vol/100
  })
  const evl4 = window.runtime.EventsOn('player_pause', (paused: boolean) => {
    isPlaying.value = !paused
  })

  onBeforeUnmount(() => {
    // unsubscribe
    evl1() 
    evl2()
    evl3()
    evl4()
  })

  GetCurrentMetadata().then((metadata: model.MusicMetadata) => {
    console.info(Date.now() / 1000)
    console.info('metadata', metadata)
    duration.value = metadata.duration
    title.value = metadata.title || 'Unknown'
    artist.value = metadata.artist || 'Unknown'
  })

})

const formatTime = (time: number): string => {
  const mins = Math.floor(time / 60)
  const secs = Math.floor(time % 60)
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
}
</script>

<style scoped>
.music-player {
  position: relative;
  width: 500px;
  height: 304px;
  border-radius: 16px;
  margin: 0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);  
  max-width: 100vw;
  max-height: 100vh;
}
.window-actions {
  position: absolute;
  top: 2px;
  right: 2px;
  z-index: 999;
}
.q-card .title {
  cursor:grab;
}

</style>