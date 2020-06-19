<script>
  export default {
    name: 'MenuItem',
    functional: true,
    props: {
      icon: {
        type: String,
        default: ''
      },
      title: {
        type: String,
        default: ''
      }
    },
    render(h, context) {
      const { icon, title } = context.props
      const vnodes = []

      if (icon) {
        // vnodes.push(<svg-icon icon-class={icon}/>)
        const style = {
          'margin-right': '5px',
          'z-index': 999
        }
        vnodes.push(<span class={icon} style={style}/>)
      }

      if (title) {
        vnodes.push(<span class='title' slot='title'>{(title)}</span>)
      }
      return vnodes
    }
  }
</script>
