import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import Misc from '../../misc';

import tags_setup_pic from '../../../assets/tags_setup_example.png';

const Picture = Misc.Picture;

class TagsSetupHelpComponent extends Component {
    constructor(props) {
        super(props);
        this.state = { open: false };
    }

    openDialog = (e) => {
        if (e)
            e.preventDefault();
        this.setState({open: true});
    };

    closeDialog = (e=null) => {
        if (e)
            e.preventDefault();
        this.setState({open: false});
    };
    
    render() {
        return (
            <div className="tags-setup-help inline-block m-l-20">
                <a href="" onClick={this.openDialog}>
                    Don't see your tags ? Click here.
                </a>
        
                <Dialog open={this.state.open} fullWidth>
                    <DialogTitle disableTypography><h1>Tags setup</h1></DialogTitle>
                    <DialogContent>
                        <div className="tutorial">
                            <ol>
                                <li>Go to your <a rel="noopener noreferrer" target="_blank" href="https://console.aws.amazon.com/billing/home#/preferences/tags">AWS Console Billing Cost Allocation Tags page.</a>.</li>
                                <li>
                                    In the table select the tag keys you want to see in TrackIt and click Activate.
                                    <Picture
                                        src={tags_setup_pic}
                                        alt="Tags setup tutorial"
                                        button={<strong>( Click here to see screenshot )</strong>}
                                    />
                                </li>
                                <li>
                                    You can click the <strong>Refresh</strong> button if you want AWS to enable Tags export quicker (you can only do this once every 24h).
                                </li>
                            </ol>
                            <div className="alert alert-info">
                                <i className="fa fa-info-circle"></i>
                                &nbsp;
                                Please note that after enabling tag export it can take up to 24h before tags data appear in TrackIt.
                            </div>
                        </div>

                        <hr/>
                        <button className="btn btn-default pull-right" onClick={this.closeDialog}>Close</button>
                    </DialogContent>
                </Dialog>
    
            </div>
    
        );
    }
}

export default TagsSetupHelpComponent;