<template>
    <div class="container col-md-7 col-sm-12 sm-container">

      <Header/>

      <div v-if="loadingGroups || loadingServices || loadingMessages" class="row mt-5 mb-5">
        <div class="col-12 mt-5 mb-2 text-center">
          <font-awesome-icon icon="circle-notch" class="text-dim" size="2x" spin/>
        </div>
        <div class="col-12 text-center mt-3 mb-3">
          <span class="text-dim">{{ loadingMessage }}</span>
        </div>
      </div>

      <div v-else-if="groups.length === 0 && services.length === 0 && messages === null" class="row mt-5 mb-5">
        <div class="col-12 text-center mt-3 mb-3">
          <span class="text-dim">No group, service or message to display.</span>
        </div>
      </div>

      <div v-else>
        <div class="col-12 full-col-12">
          <MessageBlock v-for="message in messages" v-bind:key="message.id" :message="message" />
        </div>

        <div class="col-12 full-col-12" v-if="services_no_group.length > 0">
          <div v-for="service in services_no_group" v-bind:key="service.id" class="list-group online_list mb-4">
              <div class="list-group-item list-group-item-action">
                  <router-link class="no-decoration font-3" :to="serviceLink(service)">
                    {{service.name}}
                    <MessagesIcon :messages="service.messages"/>
                  </router-link>
                  <span class="badge float-right" :class="{'bg-success': service.online, 'bg-danger': !service.online }">{{service.online ? "ONLINE" : "OFFLINE"}}</span>
                  <GroupServiceFailures :service="service"/>
                  <IncidentsBlock :service="service"/>
              </div>
          </div>
      </div>

      <Group v-for="group in groups" v-bind:key="group.id" :group="group" />

      <div class="col-12 full-col-12">
          <div v-for="service in services" :ref="service.id" v-bind:key="service.id">
              <ServiceBlock :service="service" />
          </div>
        </div>
      </div>
    </div>
</template>

<script>
import Api from "@/API";

const Group = () => import(/* webpackChunkName: "index" */ '@/components/Index/Group')
const Header = () => import(/* webpackChunkName: "index" */ '@/components/Index/Header')
const MessageBlock = () => import(/* webpackChunkName: "index" */ '@/components/Index/MessageBlock')
const ServiceBlock = () => import(/* webpackChunkName: "index" */ '@/components/Service/ServiceBlock')
const GroupServiceFailures = () => import(/* webpackChunkName: "index" */ '@/components/Index/GroupServiceFailures')
const IncidentsBlock = () => import(/* webpackChunkName: "index" */ '@/components/Index/IncidentsBlock')
const MessagesIcon = () => import(/* webpackChunkName: "index" */ '@/components/Index/MessagesIcon')

export default {
    name: 'Index',
    components: {
      IncidentsBlock,
      GroupServiceFailures,
      ServiceBlock,
      MessageBlock,
      MessagesIcon,
      Group,
      Header
    },
    data() {
        return {
            loadingGroups: true,
            loadingServices: true,
            loadingMessages: true,
            messages: null // Initialize messages to null
        };
    },
    computed: {
      loadingMessage() {
        if (this.loadingGroups) {
          return "Loading Groups";
        } else if (this.loadingServices) {
          return "Loading Services";
        } else if (this.loadingMessages) {
          return "Loading Announcements";
        }
        return ""; // To avoid an error if no loading message is displayed
      },
      groups() {
        return this.$store.getters.groupsInOrder;
      },
      services() {
        return this.$store.getters.servicesInOrder;
      },
      services_no_group() {
        return this.$store.getters.servicesNoGroup
      },
      core() {
        return this.$store.getters.core
      },
    },
    async mounted() {
      try {
        await this.$store.dispatch('loadGroups');
      } catch (error) {
        console.error("Erreur lors du chargement des groupes :", error);
      } finally {
        this.loadingGroups = false;
      }

      try {
        await this.$store.dispatch('loadServices');
      } catch (error) {
        console.error("Erreur lors du chargement des services :", error);
      } finally {
        this.loadingServices = false;
      }

      try {
        await this.$store.dispatch('loadMessages'); // Dispatch to load messages
        this.messages = this.$store.getters.messages.filter(m => this.inRange(m) && m.service === 0);
      } catch (error) {
        console.error("Erreur lors du chargement des messages :", error);
      } finally {
        this.loadingMessages = false;
      }
    },
    methods: {
      serviceLink(service) {
        return `/services/${service.id}`
      },
      inRange(message) {
        return this.isBetween(this.now(), message.start_on, message.start_on === message.end_on ? this.maxDate().toISOString() : message.end_on)
      },
      now() {
        return new Date();
      },
      maxDate() {
        return new Date(8640000000000000);
      },
      isBetween(value, min, max) {
        return value >= new Date(min) && value <= new Date(max);
      }
    }
}
</script>