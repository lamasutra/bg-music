<template>
    <q-card class="q-mx-auto options-panel no-shadow" flat square>
      <q-card-section>
        <q-select
          v-model="selectedDevice"
          :options="inputDevices"
          label="Input Device"
          option-label="name"
          option-value="path"
          outlined
          dense
        />
  
        <q-input
          v-model="shortcuts.pause"
          name="pause"
          label="Pause / Play"
          outlined
          dense
          class="q-mt-md"
          @focus="onInputFocus"
          @blur="onInputBlur"
          @keydown.prevent
        />
        <q-input
          v-model="shortcuts.mute"
          name="mute"
          label="Mute / Unmute"
          outlined
          dense
          class="q-mt-md"
          @focus="onInputFocus"
          @blur="onInputBlur"
          @keydown.prevent
        />
        <q-input
          v-model="shortcuts.prev"
          name="prev"
          label="Prev"
          outlined
          dense
          class="q-mt-sm"
          @focus="onInputFocus"
          @blur="onInputBlur"
          @keydown.prevent
        />
        <q-input
          v-model="shortcuts.next"
          name="next"
          label="Next"
          outlined
          dense
          class="q-mt-sm"
          @focus="onInputFocus"
          @blur="onInputBlur"
          @keydown.prevent
        />
        <q-input
          v-model="shortcuts['volume-up']"
          name="volume-up"
          label="Volume Up"
          outlined
          dense
          class="q-mt-sm"
          @focus="onInputFocus"
          @blur="onInputBlur"
          @keydown.prevent
        />
        <q-input
          v-model="shortcuts['volume-down']"
          name="volume-down"
          label="Volume Down"
          outlined
          dense
          class="q-mt-sm"
          @focus="onInputFocus"
          @blur="onInputBlur"
          @keydown.prevent
        />
      </q-card-section>
  
      <q-card-actions align="right">
        <q-btn label="Close" flat color="default" @click="close" />
        <q-btn label="Save" flat color="primary" @click="saveSettings" />
      </q-card-actions>
    </q-card>
  </template>
  
  <script lang="ts" setup>
  import { ref, watch, onMounted, onUnmounted } from 'vue'
  import { GetInputDevices } from '../../wailsjs/go/input/InputManagerWrapper'
  import { GetInputDeviceControls } from '../../wailsjs/go/model/Config'
  import { SaveInputControls } from '../../wailsjs/go/app/AppState'
  import type { input, model } from '../../wailsjs/go/models'
  import { EventsEmit } from '../../wailsjs/runtime/runtime'
  
  const inputDevices = ref<input.InputDevice[]>([])
  const selectedDevice = ref<input.InputDevice|null>(null)
  const shortcuts = ref<{ [key: string]: string }>({})

  let combo = ''
  let currentKeyInput: HTMLInputElement | null = null
  
  onMounted(async () => {
    inputDevices.value = await GetInputDevices()
    const settings: model.InputDeviceControls = await GetInputDeviceControls()
    const device: model.InputDevice = settings.device
    selectedDevice.value = device
    // shortcuts.value = settings.shortcuts
    shortcuts.value = flipMapping(settings.mapping)

    const evl1 = window.runtime.EventsOn('input_pressed', (data: string) => {
      console.info('input pressed', data)
      setComboIfApplicable(data)
    })
    onUnmounted(() => {
      evl1()
    })
  })

  function flipMapping(values: {[key: string]: string}): { [key: string]: string } {
    const flipped: { [key: string]: string } = {}
    for (const key in values) {
      var value: string = values[key]
      flipped[value.toString()] = key
    }
    return flipped
  }

  function resetCombo() {
    combo = ''
  }

  function setComboIfApplicable(cmb: string) {
    const currentLength = combo.split('+').length
    const length = cmb.split('+').length
    if (currentLength <= length) {
      combo = cmb
      if (currentKeyInput && combo.length) {
        shortcuts.value[currentKeyInput.name] = combo
      }
    }
  }

  watch(selectedDevice, (newVal: input.InputDevice | null) => {
    // console.info('selected device changed', oldVal, newVal)
    console.info(newVal)
  })
  
  const saveSettings = async () => {
    const payload: model.InputDeviceControls = {
      device: {
        name: selectedDevice.value?.name || '',
        path: selectedDevice.value?.path || ''
      },
      mapping: flipMapping(shortcuts.value)
    }

    console.info(payload)

    await SaveInputControls(payload)
    // Optionally show a snackbar or toast here
    await close()
  }

  const close = async () => {
    EventsEmit('show-options', false)
  }

  const onInputFocus = (event: Event) => {
    resetCombo()
    currentKeyInput = event.target as HTMLInputElement
    console.info(shortcuts.value)
  }

  const onInputBlur = (event: Event) => {
    // resetCombo()
    // currentKeyInput = null
    // console.info(shortcuts.value)
  }

  </script>
  
  <style scoped>
.options-panel {
  position: relative;
  width: 500px;
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