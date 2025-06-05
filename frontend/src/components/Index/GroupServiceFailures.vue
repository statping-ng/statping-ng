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
              <span class="font-2 text-muted">90 {{$t('days_ago')}}</span>
              <span class="divider"></span>
              <span class="text-center font-2" :class="textClass(service)">{{service_txt}}</span>
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
      this.hover_text = `${e.date.toLocaleDateString()} - ${txt}`
    },
      // These two functions have been got in PR #1234 realised by @OstlerDev
      async lastDaysFailures() {
        const start = this.beginningOf('day', this.nowSubtract(86400 * 90))
        const end = this.endOf('tomorrow')
        // Call both endpoints to get both success and failure data for the past 90 days
        const failuresPromise = Api.service_failures_data(this.service.id, this.toUnix(start), this.toUnix(end), "24h", true);
        const hitsPromise = Api.service_hits(this.service.id, this.toUnix(start), this.toUnix(end), "24h", true);
        // Wait for both promises to resolve
        const [failuresData, hitsData] = await Promise.all([failuresPromise, hitsPromise]);
        // Merge the data
        const mergedData = this.mergeData(failuresData, hitsData);
        mergedData.forEach((d) => {
          let date = this.parseISO(d.timeframe)
          // Throw out data that is from the future (shouldn't happen, but good to check)
          if ((this.toUnix(date) * 1000) > Date.now()) { 
            return 
          }
          
          this.failureData.push({
            month: date.getMonth(),
            day: date.getDate(),
            date: date,
            amount: d.amount,
            outage_type: d.outage_type,
            hits: d.hits || 0
          })
        })
      },
      mergeData(failuresData, hitsData) {
        const dataMap = new Map();
        
        // Process hits data
        hitsData.forEach(d => {
          dataMap.set(d.timeframe, { hits: d.amount, amount: 0, date: d.timeframe });
        });
        
        // Process failures data
        failuresData.forEach(d => {
          let data = dataMap.get(d.timeframe) || { hits: 0, amount: 0, date: d.timeframe };
          data.amount = d.amount;
          data.outage_type = d.outage_type;
          dataMap.set(d.timeframe, data);
        });
        
        // Convert map to array
        return Array.from(dataMap, ([date, data]) => {
          return {
            hits: data.hits,
            amount: data.amount,
            outage_type: data.outage_type,
            date: data.date,
            timeframe: date
          };
        });
      },
      // Returns a CSS class for the bar depending on the day's failure data.
      // If an outage was recorded for that day:
      // - 'Critical' returns 'day-error'
      // - 'Minor' or 'Major' returns 'day-outage'
      // Otherwise, it returns 'day-error' if there are failures, or 'day-success' if not.
      getDayClass(dayData) {
        // No data points for day
        if (dayData.amount === 0 && dayData.hits === 0) {
          return 'day-no-data';
        } 
        // No failures for day
        else if (dayData.amount === 0 && dayData.hits > 0) {
          return 'day-success';
        } 
        // Some failures for the day
        else {
          // dayData.outage_type might be 'critical', 'major', 'minor', or ''.
          if (dayData.outage_type === 'critical') {
            return 'day-critical-outage';
          } else if (dayData.outage_type === 'major') {
            return 'day-major-outage';
          } else if (dayData.outage_type === 'minor') {
            return 'day-minor-outage';
          }
          return dayData.amount > 0 ? 'day-error' : 'day-success';
        }
      }
    }
}
</script>
