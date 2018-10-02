import React, {Component} from 'react';
import Dialog from './Dialog';
import PropTypes from "prop-types";

class DeleteConfirmation extends Component {

  render() {
    return(
      <Dialog
        buttonName={<span><i className="fa fa-times"/>&nbsp;Delete</span>}
        buttonType="danger"
        title={"Are you sure to want to delete " + (this.props.entity || "this") + " ?"}
        actionName="Confirm"
        actionFunction={this.props.confirm}
        disabled={this.props.disabled}
      />
    );
  }

}

DeleteConfirmation.propTypes = {
  entity: PropTypes.string,
  confirm: PropTypes.func.isRequired,
  disabled: PropTypes.bool
};

DeleteConfirmation.defaultProps = {
  disabled: false
};

export default DeleteConfirmation;
