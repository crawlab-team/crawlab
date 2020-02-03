import i18n from '../i18n'
import store from '../store'
import stats from './stats'

export default {
  isFinishedTour: (tourName) => {
    if (!localStorage.getItem('tour')) {
      localStorage.setItem('tour', JSON.stringify({}))
      return false
    }

    let data
    try {
      data = JSON.parse(localStorage.getItem('tour'))
    } catch (e) {
      localStorage.setItem('tour', JSON.stringify({}))
      return false
    }
    return !!data[tourName]
  },
  finishTour: (tourName) => {
    let data
    try {
      data = JSON.parse(localStorage.getItem('tour'))
    } catch (e) {
      localStorage.setItem('tour', JSON.stringify({}))
      data = {}
    }
    data[tourName] = 1
    localStorage.setItem('tour', JSON.stringify(data))

    // 发送统计数据
    const finalStep = store.state.tour.tourFinishSteps[tourName]
    const currentStep = store.state.tour.tourSteps[tourName]
    if (currentStep === finalStep) {
      stats.sendEv('教程', '完成', tourName)
    } else {
      stats.sendEv('教程', '跳过', tourName)
    }
  },
  nextStep: (tourName, currentStep) => {
    store.commit('tour/SET_TOUR_STEP', {
      tourName,
      step: currentStep + 1
    })
    stats.sendEv('教程', '下一步', tourName)
  },
  prevStep: (tourName, currentStep) => {
    store.commit('tour/SET_TOUR_STEP', {
      tourName,
      step: currentStep - 1
    })
    stats.sendEv('教程', '上一步', tourName)
  },
  getOptions: (isShowHighlight) => {
    return {
      labels: {
        buttonSkip: i18n.t('Skip'),
        buttonPrevious: i18n.t('Previous'),
        buttonNext: i18n.t('Next'),
        buttonStop: i18n.t('Finish')
      },
      highlight: isShowHighlight
    }
  }
}
