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
          data-test-id="button-new-application"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.application.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::system:application/*"
          button-variant="light"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>
      </span>

      <b-dropdown
        v-if="false"
        variant="link"
        right
        menu-class="shadow-sm"
        :text="$t('export')"
      >
        <b-dropdown-item-button variant="link">
          {{ $t('yaml') }}
        </b-dropdown-item-button>
      </b-dropdown>
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
        searchPlaceholder: $t('filterForm.query.placeholder'),
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
      class="custom-resource-list-height application-list"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          data-test-id="filter-deleted-apps"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: a }">
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
            <c-permissions-button
              v-if="a.applicationID && canGrant"
              :title="a.name || a.applicationID"
              :target="a.name || a.applicationID"
              :resource="`corteza::system:application/${a.applicationID}`"
              button-variant="link text-decoration-none text-dark regular-font rounded-0"
              class="text-dark d-print-none border-0"
            >
              <font-awesome-icon :icon="['fas', 'lock']" />
              {{ $t('permissions') }}
            </c-permissions-button>
          </b-dropdown-item>

          <b-dropdown-item>
            <c-input-confirm
              borderless
              variant="link"
              size="md"
              button-class="text-decoration-none text-dark regular-font rounded-0"
              class="w-100"
              @confirmed="handleDelete(a)"
            >
              <font-awesome-icon
                :icon="['far', 'trash-alt']"
                class="text-danger"
              />
              <span
                v-if="!a.deletedAt"
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
    namespaces: 'system.applications',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'applications',

      primaryKey: 'applicationID',
      editRoute: 'system.application.edit',

      filter: {
        query: '',
        deleted: 0,
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      fields: [
        {
          key: 'name',
          sortable: true,
        },
        {
          key: 'unify.name',
          label: this.$t(`columns.appListName`),
        },
        {
          key: 'enabled',
          formatter: (v) => v ? 'Yes' : 'No',
        },
        {
          key: 'createdAt',
          sortable: true,
          formatter: (v) => moment(v).fromNow(),
        },
        {
          key: 'actions',
          label: '',
          class: 'text-right ',
        },
      ].map(c => ({
        ...c,
        // Generate column label translation key
        label: c.label || this.$t(`columns.${c.key}`),
      })),
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'application.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$SystemAPI.applicationList(this.encodeListParams()))
    },

    getAppInfo (item) {
      return { applicationID: item[this.primaryKey], name: item.name }
    },

    handleDelete (application) {
      this.incLoader()
      const { deletedAt = '' } = application
      const method = deletedAt ? 'applicationUndelete' : 'applicationDelete'
      const event = deletedAt ? 'undelete' : 'delete'
      const { applicationID } = application

      this.$SystemAPI[method]({ applicationID })
        .then(() => {
          this.toastSuccess(this.$t(`notification:application.${event}.success`))
          this.$refs.resourceList.refresh()
        })
        .catch(this.toastErrorHandler(this.$t(`notification:application.${event}.error`)))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>

<style lang="scss">
.application-list {
  td:nth-of-type(5) {
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

  tr:hover td:nth-of-type(5) {
    opacity: 1;
    background-color: $gray-200;
  }
}
</style>
