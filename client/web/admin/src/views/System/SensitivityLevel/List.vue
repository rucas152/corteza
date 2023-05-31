<template>
  <b-container
    fluid="xl"
  >
    <c-content-header
      :title="$t('title')"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="canCreate"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.sensitivityLevel.new' }"
        >
          {{ $t('new') }}
        </b-button>
      </span>
    </c-content-header>

    <c-resource-list
      ref="resourceList"
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="fields"
      :items="items"
      :row-class="genericRowClass"
      :translations="{
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
      }"
      clickable
      sticky-header
      hide-search
      class="custom-resource-list-height sensitivityLevel-list"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: s }">
        <b-dropdown
          variant="outline-light"
          toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2"
          no-caret
          dropleft
          lazy
          menu-class="m-0"
        >
          <template #button-content>
            <font-awesome-icon
              :icon="['fas', 'ellipsis-v']"
            />
          </template>

          <b-dropdown-item>
            <c-input-confirm
              borderless
              variant="link"
              size="md"
              button-class="text-decoration-none text-dark regular-font rounded-0"
              class="w-100"
              @confirmed="handleDelete(s)"
            >
              <font-awesome-icon
                :icon="['far', 'trash-alt']"
                class="text-danger"
              />
              <span
                v-if="!s.deletedAt"
                class="p-1"
              >{{ $t('delete') }}</span>

              <span
                v-else
                class="p-1"
              >{{ $t('undelete') }}</span>
            </c-input-confirm>
          </b-dropdown-item>
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import * as moment from 'moment'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import { mapGetters } from 'vuex'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.sensitivityLevel',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'sensitivityLevel',

      primaryKey: 'sensitivityLevelID',
      editRoute: 'system.sensitivityLevel.edit',

      filter: {
        query: '',
        deleted: 0,
      },

      sorting: {
        sortBy: 'level',
        sortDesc: true,
      },

      fields: [
        {
          key: 'meta.name',
        },
        {
          key: 'level',
          sortable: true,
        },
        {
          key: 'createdAt',
          sortable: true,
          formatter: (v) => moment(v).fromNow(),
        },
        {
          key: 'actions',
        },
      ].map(c => ({
        ...c,
        // Generate column label translation key
        label: this.$t(`columns.${c.key}`),
      })),
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'dal-sensitivity-level.manage')
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$SystemAPI.dalSensitivityLevelList(this.encodeListParams()))
    },

    handleDelete (sensitivityLevel) {
      this.incLoader()
      const { deletedAt = '' } = sensitivityLevel
      const method = deletedAt ? 'dalSensitivityLevelUndelete' : 'dalSensitivityLevelDelete'
      const event = deletedAt ? 'undeleted' : 'deleted'
      const { sensitivityLevelID } = sensitivityLevel

      this.$SystemAPI[method]({ sensitivityLevelID })
        .then(() => {
          this.toastSuccess(this.$t(`notification:sensitivityLevel.${event}.success`))
          this.$refs.resourceList.refresh()
        })
        .catch(this.toastErrorHandler(this.$t(`notification:sensitivityLevel.${event}.error`)))
        .finally(() => this.decLoader())
    },
  },
}
</script>

<style lang="scss">
.sensitivityLevel-list {
  td:nth-of-type(4) {
    padding-top: 8px;
    position: sticky;
    right: 0;
    opacity: 0;
    transition: opacity 0.25s;
    width: 1%;

    .regular-font {
      font-family: $font-regular !important;
    }
  }

  tr:hover td:nth-of-type(4) {
    opacity: 1;
    background-color: $gray-200;
  }
}
</style>
