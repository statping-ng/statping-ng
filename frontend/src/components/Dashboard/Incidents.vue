<template>
    <div class="col-12">

        <div v-for="incident in incidents" :key="incident.id" class="card contain-card mb-4">
            <div class="card-header">Incident: {{incident.title}}
                <button @click="deleteIncident(incident)" class="btn btn-sm btn-danger float-right">
                    <font-awesome-icon icon="times" />
                </button>
            </div>

            <FormIncidentUpdates :incident="incident"/>

            <span class="font-2 p-2 pl-3">Created: {{niceDate(incident.created_at)}} | Last Update: {{niceDate(incident.updated_at)}}</span>
        </div>


        <div class="card contain-card">
            <div class="card-header">Create Incident</div>
            <div class="card-body">
                <form @submit.prevent="createIncident">
                    <div class="form-group row">
                        <label class="col-sm-4 col-form-label">Title</label>
                        <div class="col-sm-8">
                            <input v-model="incident.title" type="text" name="title" class="form-control" id="title" placeholder="Incident Title" required>
                        </div>
                    </div>

                    <div class="form-group row">
                        <label class="col-sm-4 col-form-label">Description</label>
                        <div class="col-sm-8">
                            <textarea v-model="incident.description" rows="5" name="description" class="form-control" id="description" required></textarea>
                        </div>
                    </div>

                    <div class="form-group row">
                        <div class="col-sm-12">
                            <button :disabled="submitting || !canCreateIncident"
                                    type="submit" class="btn btn-block btn-primary">
                                {{ submitting ? "Creating Incident..." : "Create Incident" }}
                            </button>
                        </div>
                    </div>
                    <div v-if="errorMessage" class="alert alert-danger" role="alert">{{ errorMessage }}</div>
                </form>
            </div>
        </div>

    </div>
</template>

<script>
import Api from "../../API";

const FormIncidentUpdates = () => import(/* webpackChunkName: "dashboard" */ '@/forms/IncidentUpdates')

    export default {
        name: 'Incidents',
        components: {FormIncidentUpdates},
        data() {
            return {
                serviceID: 0,
                submitting: false,
                errorMessage: "",
                incidents: [],
                incident: {
                    title: "",
                    description: "",
                    service: 0,
                  }
              }
          },
    computed: {
        canCreateIncident() {
            return this.incident.title.trim().length > 0 && this.incident.description.trim().length > 0
        }
    },

    created() {
        this.serviceID = Number(this.$route.params.id);
        this.incident.service = Number(this.$route.params.id);
    },

    async mounted() {
        await this.loadIncidents()
    },

    methods: {
      extractErrorMessage(error, fallback) {
        const responseData = error?.response?.data || error
        if (typeof responseData === "string" && responseData.trim()) {
          return responseData.trim()
        }
        if (typeof responseData?.error === "string" && responseData.error.trim()) {
          return responseData.error
        }
        if (responseData?.error?.message) {
          return responseData.error.message
        }
        if (responseData?.message) {
          return responseData.message
        }
        if (error?.message) {
          return error.message
        }
        return fallback
      },

      async delete(i) {
        this.res = await Api.incident_delete(i)
        if (this.res.status === "success") {
          this.incidents = this.incidents.filter(obj => obj.id !== i.id);
          //await this.loadIncidents()
        }
      },
        async deleteIncident(incident) {
          const modal = {
            visible: true,
            title: "Delete Incident",
            body: `Are you sure you want to delete Incident ${incident.title}?`,
            btnColor: "btn-danger",
            btnText: "Delete Incident",
            func: () => this.delete(incident),
          }
          this.$store.commit("setModal", modal)
        },

        async createIncident() {
            if (this.submitting) {
                return
            }

            const title = this.incident.title.trim()
            const description = this.incident.description.trim()

            if (!title || !description) {
                this.errorMessage = "Incident title and description are required."
                return
            }

            this.submitting = true
            this.errorMessage = ""

            try {
                const response = await Api.incident_create(this.serviceID, {
                    ...this.incident,
                    title,
                    description,
                    service: this.serviceID,
                })

                if (response?.status === "success" && response.output) {
                    this.incidents.push(response.output)
                    this.incident = {
                        title: "",
                        description: "",
                        service: this.serviceID,
                    }
                    return
                }

                this.errorMessage = this.extractErrorMessage(response, "Unable to create the incident right now.")
            } catch (error) {
                this.errorMessage = this.extractErrorMessage(error, "Unable to create the incident right now.")
            } finally {
                this.submitting = false
            }
        },

        async loadIncidents() {
            this.incidents = await Api.incidents_service(this.serviceID)
        }

    }
}
</script>
