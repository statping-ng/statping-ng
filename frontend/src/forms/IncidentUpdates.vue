<template>
    <div class="card-body pt-3">

        <div v-if="updates.length===0" class="alert alert-link text-danger">
            No updates found, create a new Incident Update below.
        </div>

        <div v-for="update in updates" :key="update.id">
            <IncidentUpdate :update="update" :onUpdate="loadUpdates" :admin="true"/>
        </div>

        <form class="row" @submit.prevent="createIncidentUpdate">
            <div class="col-12 col-md-3 mb-3 mb-md-0">
                <select v-model="incident_update.type" class="form-control">
                    <option value="Investigating">Investigating</option>
                    <option value="Update">Update</option>
                    <option value="Unknown">Unknown</option>
                    <option value="Resolved">Resolved</option>
                </select>
            </div>
            <div class="col-12 col-md-7 mb-3 mb-md-0">
                <input v-model="incident_update.message" name="description" class="form-control" id="message" required>
            </div>

            <div class="col-12 col-md-2">
                <button :disabled="submitting || !canSubmit"
                        type="submit" class="btn btn-block btn-primary">
                    {{ submitting ? "Adding..." : "Add" }}
                </button>
            </div>
        </form>

        <div v-if="errorMessage" class="alert alert-danger mt-3 mb-0" role="alert">
            {{ errorMessage }}
        </div>

    </div>
</template>

<script>
    import Api from "../API";
    const IncidentUpdate = () => import(/* webpackChunkName: "index" */ "@/components/Elements/IncidentUpdate");

    export default {
        name: 'FormIncidentUpdates',
        components: {IncidentUpdate},
        props: {
            incident: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                updates: [],
                submitting: false,
                errorMessage: "",
                incident_update: {
                    incident: this.incident.id,
                    message: "",
                    type: "Investigating"
                }
            }
        },
        computed: {
            canSubmit() {
                return this.incident_update.message.trim().length > 0
            }
        },

        async mounted() {
            await this.loadUpdates()
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
            async createIncidentUpdate() {
                if (this.submitting) {
                    return
                }

                const message = this.incident_update.message.trim()
                if (!message) {
                    this.errorMessage = "Incident update message is required."
                    return
                }

                this.submitting = true
                this.errorMessage = ""

                try {
                    const response = await Api.incident_update_create({
                        ...this.incident_update,
                        message,
                    })

                    if (response?.status === "success" && response.output) {
                        this.updates.push(response.output)
                        this.incident_update = {
                            incident: this.incident.id,
                            message: "",
                            type: "Investigating"
                        }
                        return
                    }

                    this.errorMessage = this.extractErrorMessage(response, "Unable to add the incident update right now.")
                } catch (error) {
                    this.errorMessage = this.extractErrorMessage(error, "Unable to add the incident update right now.")
                } finally {
                    this.submitting = false
                }
            },

            async loadUpdates() {
                this.updates = await Api.incident_updates(this.incident)
            }
        }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
