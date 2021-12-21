<template>
  <div>
    <button @click="callServerMethod">Request Channel</button>
    <SentReceivedViewer
        :sent="sent"
        :received="received"
    />
  </div>
</template>

<script>
import SentReceivedViewer from "@/components/SentReceivedViewer.vue";
import {Flowable} from "rsocket-flowable";

export default {
  props: {
    socket: {
      type: Object,
      default: null
    }
  },
  components: {
    SentReceivedViewer
  },
  data() {
    return {
      sent: [],
      received: []
    };
  },
  methods: {
    callServerMethod() {
      console.log("request - stream call...");
      if (this.socket) {
        const flowablePayload = new Flowable(subscriber => {
          subscriber.onSubscribe({
            cancel: () => {/* no-op */
            },
            request: n => {

              for (let index = 0; index < n; index++) {
                let message = {};
                if (index === 0) {
                  console.log(0)
                  message = {
                    data: {
                      trace: '123',
                      method: '/hello/hellosrv/saychannel',
                      service: 'api.hello',
                      version: 'v1.0.0'
                    }
                  };
                } else {
                  console.log("send data after meta")
                  message = {
                    data: {ping: "ping data #" + index}
                  };
                }
                console.log(1)
                subscriber.onNext(message)

                if (index > 10) {
                  subscriber.onComplete();
                  break
                }
              }
            }
          });
        });
        // test flowable payload
        // flowablePayload.subscribe({
        //   onComplete: () => console.log("done"),
        //   onError: error => console.error(error),
        //   onNext: value => {
        //     console.log("got onNext value ");
        //     console.log(value);
        //   },
        //   // Nothing happens until `request(n)` is called
        //   onSubscribe: sub => {
        //     sub.request(5);
        //   }
        // });
        this.socket
            .requestChannel(flowablePayload)
            .subscribe({
              onComplete: () => {
                console.log(3)
                console.log("requestChannel onComplete");
                this.received.push("requestChannel onComplete");
              },
              onError: error => {
                console.log(4)
                console.log("got error with requestChannel");
                console.error(error);
              },
              onNext: value => {
                console.log(5)
                console.log("got next value in requestChannel..", value.data);
                this.received.push(value.data);
              },
              // Nothing happens until `request(n)` is called
              onSubscribe: sub => {
                console.log(6)
                console.log("subscribe request Channel!");
                sub.request(2);
                this.sent.push("requestChannel invoke success!");
              }
            });
      } else {
        console.log("not connected...");
      }
    }
  }
};
</script>

<style lang="scss" scoped>
</style>