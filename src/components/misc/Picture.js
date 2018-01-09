import React, {Component} from 'react';
import PropTypes from 'prop-types';
import Dialog, {DialogContent} from 'material-ui/Dialog';

class Picture extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false});
  };

  render() {

    let image = <img className="white-box-title-icon" src={this.props.src} alt={this.props.alt}/>;

    let button = (this.props.preview ? image : (this.props.button || <button className="btn btn-default">{this.props.text}</button>));

    return(
      <div className="picture-preview">

        <div onClick={this.openDialog}>
          {button}
        </div>

        <Dialog open={this.state.open} fullWidth maxWidth={false}>

          <DialogContent className="picture" onClick={this.closeDialog}>
            {image}
          </DialogContent>

        </Dialog>

      </div>
    );
  }

}

Picture.propTypes = {
  src: PropTypes.string.isRequired,
  alt: PropTypes.string,
  text: PropTypes.string,
  preview: PropTypes.bool,
  button: PropTypes.node
};

Picture.defaultProps = {
  text: "See picture",
  alt: "",
  preview: false,
  button: ""
};

export default Picture;
