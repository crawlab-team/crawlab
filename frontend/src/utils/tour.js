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
  }
}
