<template>
    <div>
      <div v-observe-visibility="{callback: visibleChart, once: true}" v-if="!loaded" class="row">
        <div class="col-12 text-center mt-3">
          <font-awesome-icon icon="circle-notch" class="text-dim" size="2x" spin/>
        </div>
      </div>
      <transition name="fade">
        <div v-if="loaded">
        <div class="d-flex mt-3">
            <div class="flex-fill service_day" v-for="(d, index) in failureData" @mouseover="mouseover(d)" @mouseout="mouseout" :class="getDayClass(d)">
                <span v-if="d.amount !== 0" class="d-none d-md-block text-center small"></span>
            </div>
        </div>
        <div class="row mt-2">
          <div class="col-12 no-select">
            <p class="divided">
              <span class="font-2 text-muted">{{this.days_to_show}} {{$t('days_ago')}}</span>
              <span class="divider"></span>
              <span class="text-center font-2" :class="{'text-muted': service.online, 'text-danger': !service.online}">{{service_txt}}</span>
              <span class="divider"></span>
              <span class="font-2 text-muted">{{$t('today')}}</span>
            </p>
          </div>
        </div>
      <div class="daily-failures small text-right text-dim">{{hover_text}}</div>
      </div>
  </transition>
    </div>
</template>

<script>
    import Api from '../../API';

export default {
  name: 'GroupServiceFailures',
  components: {

  },
    data() {
        return {
            failureData: [],
          hover_text: "",
          loaded: false,
          visible: false,
          days_to_show: 90,
          // Maximum number before showing the next category
          outageSeverity: {
            minor: { start: 1, end: 30 },
            moderate: { start: 30, end: 120 },
            major: { start: 120, end: 240 },
            critical: { start: 240 }
          }
        }
    },
  props: {
      service: {
          type: Object,
          required: true
      }
  },
  computed: {
    service_txt() {
      return this.smallText(this.service)
    }
  },
  mounted () {

    },
    methods: {
      visibleChart(isVisible, entry) {
        if (isVisible && !this.visible) {
          this.visible = true
          this.lastDaysFailures().then(() =>  this.loaded = true)
        }
      },
      mouseout() {
        this.hover_text = ""
      },
      mouseover(e) {
        let txt = `${e.amount} Failures`
        if (e.amount === 0) {
          txt = `No Issues`
        }
        this.hover_text = `${e.date.toUTCString().replace(" 00:00:00 GMT", "")} - ${txt}`
      },
      async lastDaysFailures() {
        const start = this.beginningOf('day', this.nowSubtract(86400 * this.days_to_show))
        const end = this.endOf('today')
        // Call both endpoints to get both success and failure data for the past 90 days
        const failuresPromise = Api.service_failures_data(this.service.id, this.toUnix(start), this.toUnix(end), "24h", true);
        const hitsPromise = Api.service_hits(this.service.id, this.toUnix(start), this.toUnix(end), "24h", true);

        // Wait for both promises to resolve
        const [failuresData, hitsData] = await Promise.all([failuresPromise, hitsPromise]);

        // Merge the data
        const mergedData = this.mergeData(failuresData, hitsData);

        mergedData.forEach((d) => {
          let date = new Date(d.timeframe);
          // Throw out data that is from the future (shouldn't happen, but good to check)
          if ((this.toUnix(date) * 1000) > Date.now()) { 
            return 
          }
          
          this.failureData.push({
            month: date.getUTCMonth() + 1,
            day: date.getUTCDate(),
            date: date,
            amount: d.amount,
            hits: d.hits || 0
          });
        });

        // Only show the last configured number of days
        this.failureData.slice(-this.days_to_show);
      },
      mergeData(failuresData, hitsData) {
        const dataMap = new Map();
        
        // Process hits data
        hitsData.forEach(d => {
          let date = this.parseISO(d.timeframe);
          dataMap.set(d.timeframe, { hits: d.amount, amount: 0, date: d.timeframe });
        });
        
        // Process failures data
        failuresData.forEach(d => {
          let date = this.parseISO(d.timeframe);

          let data = dataMap.get(d.timeframe) || { hits: 0, amount: 0, date: d.timeframe };
          data.amount = d.amount;
          dataMap.set(d.timeframe, data);
        });
        
        // Convert map to array
        return Array.from(dataMap, ([date, data]) => ({ ...data, timeframe: date }));
      },
      getDayClass(data) {
        // No data points for day
        if (data.amount === 0 && data.hits === 0) {
          return 'day-no-data';
        } 
        // No failures for day
        else if (data.amount === 0 && data.hits > 0) {
          return 'day-success';
        } 
        // Some failures for the day
        else {
          // Determine the severity and return the corresponding color class
          const outageSeverity = data.amount;
          if (outageSeverity >= this.outageSeverity.minor.start && outageSeverity < this.outageSeverity.minor.end) {
            return 'day-minor-outage'; // Light green
          } else if (outageSeverity >= this.outageSeverity.moderate.start && outageSeverity < this.outageSeverity.moderate.end) {
            return 'day-moderate-outage'; // Yellow
          } else if (outageSeverity >= this.outageSeverity.major.start && outageSeverity < this.outageSeverity.major.end) {
            return 'day-major-outage'; // Orange
          } else if (outageSeverity >= this.outageSeverity.critical.start) {
            return 'day-critical-outage'; // Red
          }
        }
      },
    }
}
</script>
