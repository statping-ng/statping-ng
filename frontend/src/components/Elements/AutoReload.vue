<template>
  <div class="col-12 full-col-12">
    <div class="row mb-3 justify-content-end">
      <div class="col-sm-2">
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <span class="input-group-text" id="validationTooltipUsernamePrepend">
              <font-awesome-icon icon="sync" class="mr-1" size="1x"/>
            </span>
          </div>
          <select class="custom-select" v-model="interval" v-on:change="onSettingsChange($event)" title="Auto Reload" style="font-size: 10px;">
            <option value="0">Don't Reload</option>
            <option value="10">10 Sec</option>
            <option value="30">30 Sec</option>
            <option value="60">1 Min</option>
            <option value="300">5 Min</option>
            <option value="600">10 Min</option>
            <option value="1800">30 Min</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "AutoReload",
  data() {
    return {
      interval: 0,
      timer: "",
    };
  },
  computed: {
    console: () => console,
    window: () => window,
    modal() {
      return this.$store.getters.modal;
    },
  },
  mounted() {
    this.initSettings();
    this.initAutoUpdate();
  },
  beforeDestroy() {
    this.cancelAutoUpdate();
  },
  methods: {
    initSettings() {
      const settings = JSON.parse(localStorage.getItem("settings"));
      if (settings) {
        this.interval = settings.interval;
      } else {
        this.updateSettings();
      }
    },
    updateSettings() {
      localStorage.setItem(
        "settings",
        JSON.stringify({ interval: this.interval })
      );
    },
    autoUpdate() {
      this.$router.go();
    },
    initAutoUpdate() {
      this.cancelAutoUpdate();
      if (this.interval > 0) {
        this.timer = setInterval(this.autoUpdate, this.interval * 1000);
      }
    },
    cancelAutoUpdate() {
      clearInterval(this.timer);
      this.timer = "";
    },
    onSettingsChange(e) {
      this.updateSettings();
      this.initAutoUpdate();
    },
  },
};
</script>