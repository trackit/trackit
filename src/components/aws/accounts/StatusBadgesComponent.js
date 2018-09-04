import React, {Component} from 'react';
import { connect } from 'react-redux';
import { OverlayTrigger, Popover } from 'react-bootstrap';
import PropTypes from 'prop-types';

const fontAwesomeCodeCheck = 'circle';
const fontAwesomeCodeMiddle = 'times-circle';
const fontAwesomeCodeTimes = 'times-circle';
const greenClass = 'green-color';
const orangeClass = 'orange-color';
const redClass = 'red-color';

const isObjectEmpty = (obj) => (obj ? (Object.keys(obj).length === 0 && obj.constructor === Object): null) ;

export class StatusBadgeComponent extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    // Determine Item warning level : 0 no warnings, 1 warning, 2 error
    getItemWarningLevel(item) {
        // Account is not a paying account
        if (!item.payer) {
            return 0;
        }
        if (item.billRepositories.length) {
            let hasError = false;
            for (let i = 0; i < item.billRepositories.length; i++) {
                const element = item.billRepositories[i];
                if (element.error.length) {
                    hasError = true;
                }
            }

            if (!hasError) { // Bills locations are working
                if  (isObjectEmpty(this.props.values)) { // But there is no data
                    return 1;
                }
                return 0;
            }
            return 2;
        }  
        // No bill locations for account
        return 2;
    }

    getItemClasses(item) {
        const classes = [
            `fa fa-${fontAwesomeCodeCheck} ${greenClass}`,
            `fa fa-${fontAwesomeCodeMiddle} ${orangeClass}`,
            `fa fa-${fontAwesomeCodeTimes} ${redClass}`
        ];
        return classes[this.getItemWarningLevel(item)];
    }

    getItemPopover(item) {
        const okText = (
            <div>
                <i className="fa fa-check green-color"/>
                &nbsp;
                This AWS account is properly set up and returns data.
            </div>
        );
        const middleText = (
            <div>
                <i className="fa fa-times orange-color"/>
                &nbsp;
                This AWS account is properly set up but <strong>TrackIt</strong> does not have data for the selected Timerange.
                <hr />
                Please select another timerange or wait for data to be imported.
                <hr />
                For more details please check the Status of your account on the <strong>Setup page</strong>.
            </div>
        );
        const badText = (
            <div>
                <i className="fa fa-times red-color"/>
                &nbsp;
                This AWS account is not set up properly. <strong>TrackIt</strong> could not access the billing data in the specified S3 bucket.
                <hr />
                Please check the Status of your account on the <strong>Setup page</strong>.
            </div>
        );

        let text = [okText, middleText, badText];
        
        return (
            <Popover id={`popover-trigger-${item.pretty}`} title={`${item.pretty} AWS Account`}>
                {text[this.getItemWarningLevel(item)]}
            </Popover>
        );
    }

    getItemBadge(item) {          
        return(
            <span className="account-status-badge" key={item.pretty} style={this.getItemWarningLevel(item) === 0 ? {display: 'none'} : {}}>
                <OverlayTrigger
                trigger={['hover', 'focus']}
                placement="bottom"
                overlay={this.getItemPopover(item)}
                >
                    <i 
                        className={this.getItemClasses(item)}
                    />
                </OverlayTrigger>
            </span>
        );
    }

    render() {
        const { accounts, selected } = this.props;

        let badges;
        // Some accounts are selected
        if (selected.length) {
            badges = selected.map(item => this.getItemBadge(item));
        } else { // No selected accounts, displaying all accounts
            if (accounts.status && accounts.values) {
                badges = accounts.values.map(item => this.getItemBadge(item));
            }
        }

        return (
            <span className="m-l-10">
                {badges}
            </span>
        );
    }
}

StatusBadgeComponent.propTypes = {
    values: PropTypes.object.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  selected: aws.accounts.selection,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
});

export default connect(mapStateToProps, mapDispatchToProps)(StatusBadgeComponent);
