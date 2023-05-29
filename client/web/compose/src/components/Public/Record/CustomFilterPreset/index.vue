<template>
  <div>
    <b-modal
      :visible="showModal"
      :title="$t('recordList.presetFilter.addFilterPresetName')"
      body-class="p-0"
      footer-class="flex-column align-items-stretch"
      centered
      @hide="onModalHide"
    >
      <b-card
        class="pt-0"
      >
        <b-form-group
          :label="$t('recordList.presetFilter.filterName')"
        >
          <b-form-input
            v-model="filterName"
            required
          />
        </b-form-group>
      </b-card>

      <template #modal-footer>
        <div
          class="d-flex justify-content-between align-items-center"
        >
          <b-button
            variant="light"
            @click="onReset"
          >
            {{ $t('recordList.presetFilter.reset') }}
          </b-button>

          <div>
            <b-button
              variant="link"
              rounded
              class="text-primary"
              @click="onModalHide"
            >
              {{ $t('general.label.cancel') }}
            </b-button>
            <b-button
              variant="primary"
              @click="onSave"
            >
              {{ $t('general.label.save') }}
            </b-button>
          </div>
        </div>
      </template>
    </b-modal>
  </div>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'CustomFilterPreset',

  props: {
    showCustomPresetFilterModal: {
      type: Boolean,
      default: false,
    },

    customFilter: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      showModal: false,
      filterName: '',
    }
  },

  watch: {
    showCustomPresetFilterModal (val) {
      this.showModal = val
    },
  },

  created () {

  },

  methods: {
    onModalHide () {
      this.showModal = false
    },

    onReset () {
      this.filterName = ''
    },

    onSave () {
      this.$emit('save', {
        name: this.filterName,
        filter: this.customFilter,
      })

      this.onModalHide()
    },
  },

}
</script>

<style lang="scss">
.position-initial {
  position: initial;
}
</style>
