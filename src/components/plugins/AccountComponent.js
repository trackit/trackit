import React, {Component} from 'react';
import PropTypes from 'prop-types';
import Result from './ResultComponent';
  
class AccountComponent extends Component {

    getCloudScore(results) {
        let res = 0;
        for (let i = 0; i < results.length; i++) {
            const element = results[i];
            if (element.checked)
                res += ((element.passed * 100) / element.checked);
        }
        return (res / results.length);
    }

    render() {

        const { account, label, results } = this.props.account;
        const resultItems = results.map(item => <Result result={item} key={item.accountPluginIdx} />);
        const cloudScore = this.getCloudScore(results);

        return (
            <div className="conformity-account">
                <div className="white-box no-padding">
                    <div className="row conformity-account-header">
                        <div className="col-md-8">
                            <h4 className="white-box-title conformity-account-title">
                                {label ? `${label} (${account})` : `Account ${account}`}
                            </h4>
                        </div>
                        <div className="col-md-3 p-t-10">
                            <div className="progress">
                                <div
                                    className={`progress-bar blue-bg`}
                                    role="progressbar"
                                    aria-valuenow={cloudScore}
                                    aria-valuemin="0"
                                    aria-valuemax="100"
                                    style={{width: `${cloudScore}%`}}
                                >
                                </div>
                            </div>
                            <p className="no-padding no-margin">
                                Cloud score : {`${cloudScore.toFixed(0)}%`}
                            </p>
                        </div>
                    </div>

                    {resultItems}
                </div>
            </div>
        );
    }
}

AccountComponent.propTypes = {
    account: PropTypes.shape({
        account: PropTypes.string.isRequired,
        label: PropTypes.string,
        results: PropTypes.array.isRequired,
    }).isRequired,
};

export default AccountComponent;