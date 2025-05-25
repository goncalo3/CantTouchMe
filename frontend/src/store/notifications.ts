import { reactive, createApp, h } from 'vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import Alert from '@/components/Alert.vue'

export type AlertType = 'info' | 'confirm' | 'error'

export interface Alert {
  message: string
  type: AlertType
}

export interface Confirm {
  message: string
  onConfirm: () => void
  onCancel?: () => void
}

interface NotificationState {
  alert: Alert | null
  confirm: Confirm | null
  confirmVisible: boolean // Add confirmVisible to state
}

const alertFromSession = sessionStorage.getItem('alert')
const state = reactive<NotificationState>({
  alert: alertFromSession ? JSON.parse(alertFromSession) as Alert : null,
  confirm: null,
  confirmVisible: false, // Initialize confirmVisible
})

export function useNotificationStore() {
  return state
}

export function clearAlert() {
  state.alert = null
  sessionStorage.removeItem('alert')
}

export function setConfirm(message: string, onConfirm: () => void, onCancel?: () => void) {
  state.confirm = { message, onConfirm, onCancel }
}

export function clearConfirm() {
  state.confirm = null
}

// Promise-based confirm
export function showConfirm(message: string): Promise<boolean> {
  return new Promise((resolve) => {
    const container = document.createElement('div')
    document.body.appendChild(container)

    const app = createApp({
      render() {
        return h(ConfirmDialog, {
          message,
          onConfirm: () => {
            resolve(true)
            app.unmount()
            document.body.removeChild(container)
          },
          onCancel: () => {
            resolve(false)
            app.unmount()
            document.body.removeChild(container)
          },
        })
      },
    })

    app.mount(container)
  })
}

export function renderAlert(alert: Alert) {
  console.log('renderAlert called with alert:', alert)

  const container = document.createElement('div')
  document.body.appendChild(container)

  const app = createApp({
    render() {
      return h(Alert, {
        alert,
        onClose: () => {
          app.unmount()
          document.body.removeChild(container)
        },
      })
    },
  })

  app.mount(container)
}

export function showAlertWithRedirect(alert: Alert) {
  state.alert = { ...alert } // Set in the store for rendering after redirect
  sessionStorage.setItem('alert', JSON.stringify(alert))
}
