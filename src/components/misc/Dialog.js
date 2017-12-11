import React, {Component} from 'react';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogTitle,
} from 'material-ui/Dialog';
import PropTypes from "prop-types";

class DialogComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.executeAction = this.executeAction.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    if (this.props.onOpen)
      this.props.onOpen();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    if (this.props.onClose)
      this.props.onClose();
    this.setState({open: false});
  };

  executeAction = (e) => {
    e.preventDefault();
    this.closeDialog(e);
    this.props.actionFunction();
  };

  render() {
    const title = (this.props.title ? (
      <DialogTitle>
        <h1>{this.props.title}</h1>
      </DialogTitle>
    ) : null);

    const content = (this.props.children ? (
      <DialogContent>
        {this.props.children}
      </DialogContent>
    ) : null);

    const mainAction = (this.props.actionName && this.props.actionFunction ? (
      <button className="btn btn-default" onClick={this.executeAction}>
        {this.props.actionName}
      </button>
    ) : null);

    return(
      <div>

        <button className={"btn btn-" + this.props.buttonType} onClick={this.openDialog}>
          {this.props.buttonName}
        </button>

        <Dialog open={this.state.open} fullWidth>

          {title}

          {content}

          <DialogActions>

            <button className="btn btn-default" onClick={this.closeDialog}>
              {this.props.secondActionName}
            </button>

            {mainAction}

          </DialogActions>

        </Dialog>

      </div>
    );
  }

}

DialogComponent.propTypes = {
  buttonType: PropTypes.string,
  buttonName: PropTypes.string.isRequired,
  title: PropTypes.string,
  actionName: PropTypes.string,
  actionFunction: PropTypes.func,
  secondActionName: PropTypes.string,
  children: PropTypes.node,
  onOpen: PropTypes.func,
  onClose: PropTypes.func,
};

DialogComponent.defaultProps = {
  buttonType: "default",
  secondActionName: "Cancel"
};

export default DialogComponent;
