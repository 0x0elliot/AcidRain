class Observable {
  constructor() {
      this.observers = [];
  }

  subscribe(fn) {
      this.observers.push(fn);
  }

  unsubscribe(fn) {
      this.observers = this.observers.filter(observer => observer !== fn);
  }

  notify(data) {
      this.observers.forEach(observer => observer(data));
  }
}

// Create an observable for __st.cid
const cidObservable = new Observable();

// Set up the watcher
if (window.__st) {
  let currentCid = window.__st.cid;
  Object.defineProperty(window.__st, 'cid', {
      get: function() {
          return currentCid;
      },
      set: function(newValue) {
          if (newValue !== currentCid) {
              currentCid = newValue;
              cidObservable.notify(newValue);
          }
      }
  });
}

// Now you can "listen" for changes like this:
cidObservable.subscribe((newCid) => {
  if (newCid) {
    navigator.pushManager.getSubscription()
      .then(function (subscription) {
        if (subscription) {
          syncSubscriptionOnServer(subscription, { cid: newCid });
        }
      });
  }
});

// this file will eventually belong in a CDN
document.addEventListener('DOMContentLoaded', function () {
  let baseUrl = window.location.origin;

  if ('serviceWorker' in navigator && 'PushManager' in window) {
    navigator.serviceWorker.register(`${baseUrl}/apps/acidrain/public/service-worker.js`)
      .then(function (registration) {
        console.log('Service Worker registered with scope:', registration.scope);

        // Check if the user is already subscribed
        registration.pushManager.getSubscription()
          .then(async function (subscription) {
            if (subscription) {
              console.log('Already subscribed:', subscription);
              // check local storage for subscription to avoid spam
              // if (localStorage.getItem('acidRainWebPush') === '1') {
              //   console.log('Already subscribed:', subscription);
              //   return;
              // }
            } else {
              askForNotificationPermission(registration);
            }
          });
      })
      .catch(function (error) {
        console.error('Service Worker registration failed:', error);
      });
  } else {
    console.warn('Push messaging is not supported');
  }
});

function askForNotificationPermission(registration) {
  Notification.requestPermission().then(function (permission) {
    if (permission === 'granted') {
      console.log('Notification permission granted.');
      subscribeUserToPush(registration);
    } else {
      console.warn('Notification permission denied.');
    }
  });
}

function subscribeUserToPush(registration) {
  let applicationServerKey;
  let storeUrl;
  storeUrl = window.location.origin;

  // get public key from the server
  fetch(`${storeUrl}/apps/acidrain/api/web-push-public-key`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }).then(function (response) {
    if (!response.ok) {
      throw new Error('Failed to get public key from server');
    }
    return response.json();
  }).then(function (data) {
    applicationServerKey = urlB64ToUint8Array(data.publicKey);
    registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: applicationServerKey
    })
      .then(function (subscription) {
        console.log('User is subscribed:', subscription);
        sendSubscriptionToServer(subscription);
      })
      .catch(function (error) {
        console.error('Failed to subscribe the user:', error);
      });
  }).catch(function (error) {
    console.error('Error getting public key from server:', error);
  });
}

function syncSubscriptionOnServer(subscription, customer) {
  let storeUrl;
  if (window.Shopify) {
    storeUrl = window.Shopify.shop;
  } else {
    // not ideal
    storeUrl = window.location.origin;
  }

  requestObj = {
    subscription: subscription,
    storeUrl: storeUrl,
    customer: {
      cid: customer.cid
    }
  };

  let baseUrl = window.location.origin;
  fetch(`${baseUrl}/apps/acidrain/api/notification/sync`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestObj)
  })
    .then(function (response) {
      if (!response.ok) {
        throw new Error('Failed to sync subscription on server');
      }
      return response.json();
    })
    .then(function (data) {
      console.log('Subscription synced on server:', data);
    })
    .catch(function (error) {
      console.error('Error syncing subscription on server:', error);
    });

}

function urlB64ToUint8Array(base64String) {
  const padding = '='.repeat((4 - base64String.length % 4) % 4);
  const base64 = (base64String + padding)
    .replace(/-/g, '+')
    .replace(/_/g, '/');

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}

function sendSubscriptionToServer(subscription) {
  let storeUrl;
  if (window.Shopify) {
    storeUrl = window.Shopify.shop;
  } else {
    storeUrl = window.location.origin;
  }

  cid = window.__st.cid;

  requestObj = {
    subscription: subscription,
    storeUrl: storeUrl,
  };

  let baseUrl = window.location.origin;
  fetch(`${baseUrl}/apps/acidrain/api/notification/subscribe`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestObj)
  })
    .then(function (response) {
      if (!response.ok) {
        throw new Error('Failed to send subscription to server');
      }
      return response.json();
    })
    .then(function (data) {
      localStorage.setItem('acidRainWebPush', '1');
      localStorage.setItem('acidRainWebPushSubscription', JSON.stringify(subscription));
      console.log('Subscription sent to server:', data);
    })
    .catch(function (error) {
      console.error('Error sending subscription to server:', error);
    });
}
